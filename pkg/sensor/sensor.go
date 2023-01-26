package sensor

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/darox/miflora-go/pkg/error"
	"github.com/darox/miflora-go/pkg/model"
	"github.com/go-ble/ble"
)

// writes the characteristic to allow the read of temperature, light, moisture and conductivity sensor data
func enableRead(cln ble.Client, p *ble.Profile) {
	// UUID of characteristic to enable read of temperature, light, moisture and conductivity sensor data
	const u = "00001a0000001000800000805f9b34fb"
	enableReadUUID := ble.MustParse(u)
	// bytes to enable read
	enableReadBytes := []byte{0xa0, 0x1f}
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			if c.UUID.Equal(enableReadUUID) {
				if (c.Property & ble.CharWrite) != 0 {
					err := cln.WriteCharacteristic(c, enableReadBytes, false)
					error.Check(err)
				}
			}
		}
	}
}

// reads the characteristic to get firmware version and battery level
func readBattFw(cln ble.Client, p *ble.Profile) {
	// UUID of characteristic holding the battery level and firmware version
	const u = "00001a0200001000800000805f9b34fb"
	battFwUUID := ble.MustParse(u)
	for _, s := range p.Services {
		// UUID of battery level and firmware version characteristic
		for _, c := range s.Characteristics {
			if c.UUID.Equal(battFwUUID) {
				if (c.Property & ble.CharRead) != 0 {
					b, err := cln.ReadCharacteristic(c)
					error.Check(err)
					fmt.Printf("\n🔋  Battery Level: %d%%\n⚙️   Firmware: %s\n", b[0], b[2:])
				}
			}
		}
	}
}

// calls funcs to read firmware,battery, conductivity, humidity, light and temperature sensor data
func ReadAll(config *model.Configuration) {
	for _, device := range *config.Devices {
		filter := func(a ble.Advertisement) bool {
			return strings.EqualFold(a.Addr().String(), *device.Adddress)
		}
		// Scan for specified durantion, or until interrupted by user.
		if device.Alias != nil {
			fmt.Printf("\nScanning on device %s with alias %s \n", *device.Adddress, *device.Alias)
		} else {
			fmt.Printf("\nScanning on device %s \n", *device.Adddress)
		}
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *config.ScanDuration*time.Second))
		cln, err := ble.Connect(ctx, filter)

		if err != nil {
			fmt.Printf("Failed to connect to device with address: %s  and alias %s \n", *device.Adddress, *device.Alias)
			continue
		} else {
			fmt.Printf("Connected to device with address: %s  and alias %s \n", *device.Adddress, *device.Alias)

			// Make sure we had the chance to print out the message.
			done := make(chan struct{})
			// Normally, the connection is disconnected by us after our exploration.
			// However, it can be asynchronously disconnected by the remote peripheral.
			// So we wait(detect) the disconnection in the go routine.
			go func() {
				<-cln.Disconnected()
				fmt.Printf("\nDisconnected from device with address: %s \n", cln.Addr())
				close(done)
			}()

			p, err := cln.DiscoverProfile(true)
			error.Check(err)

			// Enable sensor read
			enableRead(cln, p)

			// Read battery and firmware
			readBattFw(cln, p)

			// Read sensor data
			readCondHumLighTemp(cln, p)

			err = cln.CancelConnection()
			error.Check(err)

			<-done
		}
	}
}

// reads the characteristic to get conductivity, humidity, light and temperature sensor data
func readCondHumLighTemp(cln ble.Client, p *ble.Profile) {
	// UUID of characteristic holding the sensor data for conductivity, humidity, light and temperature
	const u = "00001a0100001000800000805f9b34fb"
	condHumLighTempChar := ble.MustParse(u)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			if c.UUID.Equal(condHumLighTempChar) {
				if (c.Property & ble.CharRead) != 0 {
					b, err := cln.ReadCharacteristic(c)
					error.Check(err)
					// Temperature is stored as a 16-bit signed integer in units of 0.1°C
					var subtrahend float32 = 10.0
					temp := float32(binary.LittleEndian.Uint16(b[0:2])) / subtrahend
					light := binary.LittleEndian.Uint32(b[3:7])
					moisture := b[7]
					conductivity := binary.LittleEndian.Uint16(b[8:10])
					fmt.Printf("🌡️   Temperature: %.1f°C\n💡  Light: %d Lux\n💧  Moisture: %d%%\n⚡  Conductivity: %d\n",
						temp, light, moisture, conductivity)
				}
			}
		}
	}
}
