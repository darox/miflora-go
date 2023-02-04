package sensor

import (
	"testing"

	"github.com/darox/miflora-go/pkg/configuration"
)

/*
func TestGetReadings(t *testing.T) {

	var configPath = flag.String("config-path", "../../configuration/config.yaml", "Configuration file to use")
	flag.Parse()

	config := configuration.Create(configPath)

	d, err := dev.NewDevice(*config.Adapter)
	errorhandler.Check(err)

	ble.SetDefaultDevice(d)

	got, err := GetReadings(&config)

	if err != nil {
		errorhandler.Check(err)
	}

	if len(got.DevicesReadings) != len(*config.Devices) {
		t.Errorf("got %v want %v", len(got.DevicesReadings), len(*config.Devices))
	}
}
*/

func TestDefineIdentifierAlias(t *testing.T) {
	var alias string = "Monstera"
	var ptrAlias *string
	ptrAlias = &alias

	var address string = "422b23155c369dfee0aea210d1a9bc37"
	var ptrAddress *string
	ptrAddress = &address
	
	d := Device{
		Device: configuration.Device{
			Adddress: ptrAddress,
			Alias:   ptrAlias,
		},
		Identifier: nil,
	}

	d.defineIdentifier()

	if *d.Identifier != *d.Device.Alias {
		t.Errorf("got %v want %v", *d.Identifier, *d.Device.Alias)
	}
}

func TestDefineIdentifierAddress(t *testing.T) {
	var address string = "422b23155c369dfee0aea210d1a9bc37"
	var ptrAddress *string
	ptrAddress = &address
	
	d := Device{
		Device: configuration.Device{
			Adddress: ptrAddress,
			Alias:   nil,
		},
		Identifier: nil,
	}

	d.defineIdentifier()

	if *d.Identifier != *d.Device.Adddress {
		t.Errorf("got %v want %v", *d.Identifier, *d.Device.Adddress)
	}
}

