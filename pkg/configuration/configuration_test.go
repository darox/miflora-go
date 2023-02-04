package configuration

import (
	"flag"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	// Test uses a config file stored under /configuration/config.yaml
	var configPath = flag.String("config-path", "../../configuration/config.yaml", "Configuration file to use")
	flag.Parse()

	got := Create(configPath)

	var adapter string = "default"
	var ptrAdapter *string
	ptrAdapter = &adapter

	var scanDuration time.Duration = 5 * time.Duration(time.Nanosecond)
	var ptrScanDuration *time.Duration
	ptrScanDuration = &scanDuration

	var address string = "422b23155c369dfee0aea210d1a9bc37"
	var ptrAddress *string
	ptrAddress = &address

	var alias string = "Monstera"
	var ptrAlias *string
	ptrAlias = &alias

	var devices []Device
	devices = append(devices, Device{
		Adddress: ptrAddress,
		Alias:    ptrAlias,
	})

	want := Configuration{
		Adapter:      ptrAdapter,
		ScanDuration: ptrScanDuration,
		Devices:      &devices,
	}

	if *want.Adapter != *got.Adapter {
		t.Errorf("got %v want %v", got.Adapter, want.Adapter)
	}
	for i := range *want.Devices {
		if *(*want.Devices)[i].Adddress != *(*got.Devices)[i].Adddress {
			t.Errorf("got %v want %v", *(*got.Devices)[i].Adddress, *(*want.Devices)[i].Adddress)
		}
		if *(*want.Devices)[i].Alias != *(*got.Devices)[i].Alias {
			t.Errorf("got %v want %v", *(*got.Devices)[i].Alias, *(*want.Devices)[i].Alias)
		}
	}

	if *want.ScanDuration != *got.ScanDuration {
		t.Errorf("got %v want %v", got.ScanDuration, want.ScanDuration)
	}
}
