package sensor

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/darox/miflora-go/pkg/configuration"
	"github.com/go-ble/ble"
)

func GetReadings(c *configuration.Configuration) (devicesReadings DevicesReadings, err error) {
	// loop through devices
	for _, device := range *c.Devices {
		// check if device address follows the correct format. 
		filter := func(a ble.Advertisement) bool {
			return strings.EqualFold(a.Addr().String(), *device.Adddress)
		}
		
		d := Device{
			Device:    device,
			Identifier: nil,
		}
		// Define identifier to use for device based on alias or address
		d.defineIdentifier()

		if c.StructuredOutput == nil || !*c.StructuredOutput {
			printFormattedProgress("scanning", *d.Identifier)
		} else {
			printStructuredProgress("scanning", *d.Identifier)
		}

		// Setup context with timeout
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *c.ScanDuration*time.Second))
		// scan for device
		cln, err := ble.Connect(ctx, filter)

		if err != nil {
			if c.StructuredOutput == nil || !*c.StructuredOutput {
				printFormattedProgress("failed", *d.Identifier)
				continue
			} else {
				printStructuredProgress("failed", *d.Identifier)
				continue
			}
		} else {
			if c.StructuredOutput == nil || !*c.StructuredOutput {
				printFormattedProgress("connected", *d.Identifier)
			} else {
				printStructuredProgress("connected", *d.Identifier)
			}

			done := make(chan struct{})

			go func() {
				<-cln.Disconnected()
				if c.StructuredOutput == nil || !*c.StructuredOutput {
					printFormattedProgress("disconnected", *d.Identifier)
				} else {
					printStructuredProgress("disconnected", *d.Identifier)
				}
				close(done)
			}()

			p, err := cln.DiscoverProfile(true)
			// continue if error occurs as we wanna try to connect to the next device
			if err != nil {
				continue
			}

			err = d.enableSensorReadings(cln, p)
			// continue if error occurs as we wanna try to connect to the next device
			if err != nil {
				continue
			}

			// read battery level and firmware version
			systemReadings, err := getSystemReadings(cln, p)
			// continue if error occurs as we wanna try to connect to the next device
			if err != nil {
				continue
			}

			// Read sensor data
			sensorReadings, err := getSensorReadings(cln, p)
			if err != nil {
				continue
			}

			deviceReading := DeviceReadings{
				Alias:          *device.Alias,
				SystemReadings: systemReadings,
				SensorReadings:  sensorReadings,
			}
			// Add device readings to the list of devices readings
			devicesReadings.Add(deviceReading)

			err = cln.CancelConnection()
			if err != nil {
				continue
			}

			<-done
		}
	}
	return devicesReadings, nil
}

// Enable read of temperature, humidity, light and conductivity
func (device *Device) enableSensorReadings(cln ble.Client, p *ble.Profile) (err error) {
	// UUID service and characteristic to enable read
	cu := ble.MustParse("00001a0000001000800000805f9b34fb")
	su := ble.MustParse("0000120400001000800000805f9b34fb")
	// bytes to enable read
	enableSensorReadingsBytes := []byte{0xa0, 0x1f}

	// find service
	s := findService(p, su)
	if s == nil {
		return fmt.Errorf("service not found")
	}
	// find characteristic
	c := findCharacteristic(s, cu)
	if c == nil {
		return fmt.Errorf("characteristic not found")
	}

	// write characteristic
	err = cln.WriteCharacteristic(c, enableSensorReadingsBytes, false)
	if err != nil {
		return err
	}
	return nil
}

// Read the characteristics for battery level and firmware version
func getSystemReadings(cln ble.Client, p *ble.Profile) (systemReadings SystemReadings, err error) {
	// UUID service and characteristic to read battery level and firmware version
	cu := ble.MustParse("00001a0200001000800000805f9b34fb")
	su := ble.MustParse("0000120400001000800000805f9b34fb")

	// find service and characteristic
	s := findService(p, su)
	if s == nil {
		return systemReadings, fmt.Errorf("service not found")
	}

	c := findCharacteristic(s, cu)
	if c == nil {
		return systemReadings, fmt.Errorf("characteristic not found")
	}

	b, err := cln.ReadCharacteristic(c)
	if err != nil {
		return systemReadings, err
	}

	systemReadings = SystemReadings{
		Battery:  uint16(b[0]),
		Firmware: string(b[2:]),
	}
	return systemReadings, nil
}

// Define identifier of the device
func (device *Device) defineIdentifier() {
	if device.Device.Alias != nil {
		device.Identifier = device.Device.Alias
	} else {
		device.Identifier = device.Device.Adddress
	}
}

// Read the characteristic to get conductivity, humidity, light and temperature sensor data
func getSensorReadings(cln ble.Client, p *ble.Profile) (sensorReadings SensorReadings, err error) {
	// UUID of service and characteristic holding the sensor data
	cu := ble.MustParse("00001a0100001000800000805f9b34fb")
	su := ble.MustParse("0000120400001000800000805f9b34fb")

	// Find the service and characteristic.
	s := findService(p, su)
	if s == nil {
		return sensorReadings, fmt.Errorf("service not found")
	}

	c := findCharacteristic(s, cu)
	if c == nil {
		return sensorReadings, fmt.Errorf("characteristic not found")
	}

	if (c.Property & ble.CharRead) != 0 {
		b, err := cln.ReadCharacteristic(c)
		if err != nil {
			return sensorReadings, err
		}
		var subtrahend float32 = 10.0
		sensorReadings = SensorReadings{
			Conductivity: binary.LittleEndian.Uint16(b[8:10]),
			Moisture:     uint16(b[7]),
			Light:        binary.LittleEndian.Uint32(b[3:7]),
			Temp:         float32(binary.LittleEndian.Uint16(b[0:2])) / subtrahend,
		}
		return sensorReadings, nil
	}
	return sensorReadings, nil
}

// Print the progress of the scan in a formatted way
func printFormattedProgress(s string, d string) {
	switch s {
	case "scanning":
		fmt.Printf("📡  Scanning for %s\n", d)
	case "connected":
		fmt.Printf("✅  Connected to %s\n", d)
	case "disconnected":
		fmt.Printf("👋  Disconnected from %s\n", d)
	case "failed":
		fmt.Printf("⛔  Failed to connect to %s\n", d)
	}
}

// Print the progress of the scan in a structured way
func printStructuredProgress(s string, d string) {
	t := time.Now().Format(time.RFC3339)
	switch s {
	case "scanning":
		fmt.Printf("{\"time\":\"%s\",\"status\":\"scanning\",\"device\":\"%s\"}\n", t, d)
	case "connected":
		fmt.Printf("{\"time\":\"%s\",\"status\":\"connected\",\"device\":\"%s\"}\n", t, d)
	case "disconnected":
		fmt.Printf("{\"time\":\"%s\",\"status\":\"disconnected\",\"device\":\"%s\"}\n", t, d)
	case "failed":
		fmt.Printf("{\"time\":\"%s\",\"status\":\"failed\",\"device\":\"%s\"}\n", t, d)
	}
}

// Find the service based on the UUID
func findService(p *ble.Profile, u ble.UUID) *ble.Service {
	for _, s := range p.Services {
		if s.UUID.Equal(u) {
			return s
		}
	}
	return nil
}

// Find the characteristic based on the UUID
func findCharacteristic(s *ble.Service, u ble.UUID) *ble.Characteristic {
	for _, c := range s.Characteristics {
		if c.UUID.Equal(u) {
			return c
		}
	}
	return nil
}

// Add a new DeviceReadings to DevicesReadings
func (devicesReadings *DevicesReadings) Add(deviceReadings DeviceReadings) []DeviceReadings {
	devicesReadings.DevicesReadings = append(devicesReadings.DevicesReadings, deviceReadings)
	return devicesReadings.DevicesReadings
}

// Print the DevicesReadings in a formatted way
func (devicesReadings DevicesReadings) PrintFormatted() {
	for _, e := range devicesReadings.DevicesReadings {
		fmt.Printf("\n🪴   Name: %s \n"+
			"🔋  Battery Level: %d%% \n"+
			"⚙️   Firmware: %s \n"+
			"🌡️   Temperature: %.1f°C \n"+
			"💧  Light: %d Lux \n"+
			"⚡  Moisture: %d%% \n"+
			"🌱  Conductivity: %d µS/cm \n",
			e.Alias,
			e.SystemReadings.Battery,
			e.SystemReadings.Firmware,
			e.SensorReadings.Temp,
			e.SensorReadings.Light,
			e.SensorReadings.Moisture,
			e.SensorReadings.Conductivity)
	}
}

// Print the DevicesReadings in a structured way
func (devicesReadings DevicesReadings) PrintStructured() {
	for _, e := range devicesReadings.DevicesReadings {
		fmt.Printf("{"+
			"\"time\": \"%s\", "+
			"\"alias\": \"%s\", "+
			"\"battery\": %d, "+
			"\"firmware\": \"%s\", "+
			"\"temperature\": %.1f, "+
			"\"light\": %d, "+
			"\"moisture\": %d, "+
			"\"conductivity\": %d"+
			"}\n",
			time.Now().Format(time.RFC3339),
			e.Alias,
			e.SystemReadings.Battery,
			e.SystemReadings.Firmware,
			e.SensorReadings.Temp,
			e.SensorReadings.Light,
			e.SensorReadings.Moisture,
			e.SensorReadings.Conductivity)
	}
}

// DevicesReadings is a struct that holds the readings from all devices
type DevicesReadings struct {
	DevicesReadings []DeviceReadings
}

// DeviceReadings is a struct that holds the readings from a single device
type DeviceReadings struct {
	Alias          string
	SystemReadings SystemReadings
	SensorReadings SensorReadings
}

// Device is a struct that holds the device information
type Device struct {
	Device     configuration.Device
	Identifier *string
}

// SystemReadings is a struct that holds the battery and firmware information
type SystemReadings struct {
	Battery  uint16
	Firmware string
}

// SensorReadings is a struct that holds the conductivity, moisture, light and temperature information
type SensorReadings struct {
	Conductivity uint16
	Moisture     uint16
	Light        uint32
	Temp         float32
}
