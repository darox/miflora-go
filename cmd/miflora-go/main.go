package main

import (
	"flag"

	"github.com/darox/miflora-go/pkg/configuration"
	"github.com/darox/miflora-go/pkg/errorhandler"
	"github.com/darox/miflora-go/pkg/sensor"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

func main() {
	var configPath = flag.String("config-path", "config.yaml", "Configuration file to use")
	flag.Parse()

	config := configuration.Create(configPath)

	d, err := dev.NewDevice(*config.Adapter)
	errorhandler.Check(err)

	ble.SetDefaultDevice(d)

	r, err := sensor.GetReadings(&config)

	if err != nil {
		errorhandler.Check(err)
	}

	if config.StructuredOutput != nil && *config.StructuredOutput {
		r.PrintStructured()
	} else {
		r.PrintFormatted()
	}
}
