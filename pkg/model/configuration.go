package model

import "time"

type Configuration struct {
	Adapter      *string        `mapstructure:"adapter"`
	ScanDuration *time.Duration `mapstructure:"scanDuration"`
	Devices      *[]Device      `mapstructure:"devices"`
}

type Device struct {
	Adddress *string `mapstructure:"address"`
	Alias    *string `mapstructure:"alias"`
	Name     *string `mapstructure:"name"`
}
