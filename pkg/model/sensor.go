package model

// DevicesReadings is a struct that holds the readings from all devices
type DevicesReadings struct {
	Devices *[]DeviceReadings
}

// DeviceReadings is a struct that holds the readings from a single device
type DeviceReadings struct {
	Firmware     []byte
	Battery      byte
	Temp         float32
	Light        uint32
	Moisture     byte
	Conductivity uint16
}

