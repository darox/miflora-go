package main

import (
	"flag"

	"github.com/darox/miflora-go/pkg/error"
	"github.com/darox/miflora-go/pkg/model"
	"github.com/darox/miflora-go/pkg/sensor"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

func main() {
	var configPath = flag.String("config-path", "config.yaml", "Configuration file to use")
	flag.Parse()

	vp := viper.New()
	vp.SetConfigFile(*configPath)

	err := vp.ReadInConfig()
	error.Check(err)

	scanDurationSeconds := 5
	vp.SetDefault("scanDuration", scanDurationSeconds)
	vp.SetDefault("adapter", "default")
	vp.SetDefault("name", "Flower care")
	vp.SetDefault("address", "")

	err = vp.Unmarshal(&config)
	error.Check(err)

	validate := validator.New()
	err = validate.Struct(&config)
	error.Check(err)

	d, err := dev.NewDevice(*config.Adapter)
	error.Check(err)

	ble.SetDefaultDevice(d)

	sensor.ReadAll(&config)
}

var config model.Configuration
