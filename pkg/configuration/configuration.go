package configuration

import (
	"time"

	"github.com/darox/miflora-go/pkg/errorhandler"
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

func Create(configPath *string) (config Configuration) {
	vp := viper.New()
	vp.SetConfigFile(*configPath)

	err := vp.ReadInConfig()
	errorhandler.Check(err)

	scanDurationSeconds := 5
	vp.SetDefault("scanDuration", scanDurationSeconds)
	vp.SetDefault("adapter", "default")
	vp.SetDefault("name", "Flower care")
	vp.SetDefault("address", "")

	err = vp.Unmarshal(&config)
	errorhandler.Check(err)

	validate := validator.New()
	err = validate.Struct(&config)
	errorhandler.Check(err)

	return config
}

type Configuration struct {
	Adapter          *string        `mapstructure:"adapter" validate:"alphanumunicode"`
	Devices          *[]Device      `mapstructure:"devices" validate:"dive,required,alphanumunicode"`
	ScanDuration     *time.Duration `mapstructure:"scanDuration" validate:"numeric"`
	StructuredOutput *bool          `mapstructure:"structuredOutput"`
	Forwarder        *Forwarder     `mapstructure:"forwarder"`
}

type Forwarder struct {
	Url 	*string `mapstructure:"url" validate:"required,url"`
}

type Device struct {
	Adddress *string `mapstructure:"address"`
	Alias    *string `mapstructure:"alias"`
}
