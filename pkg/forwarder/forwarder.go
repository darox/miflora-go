package forwarder

import (
	"bytes"
	"net/http"

	"encoding/json"

	"github.com/darox/miflora-go/pkg/configuration"
	"github.com/darox/miflora-go/pkg/sensor"
)

func Post(config *configuration.Forwarder, d sensor.DevicesReadings) error {

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = http.Post(*config.Url, "application/json", bytes.NewBuffer(b))
	if err != nil {	
		return err
	}

	return nil
}
