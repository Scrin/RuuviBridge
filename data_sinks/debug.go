package data_sinks

import (
	"encoding/json"
	"fmt"

	"github.com/Scrin/RuuviBridge/parser"
	log "github.com/sirupsen/logrus"
)

func Debug() chan<- parser.Measurement {
	log.Info("Starting debug sink")
	measurements := make(chan parser.Measurement)
	go func() {
		for measurement := range measurements {
			data, err := json.Marshal(measurement)
			if err != nil {
				log.WithError(err).Error("Failed to serialize measurement")
			} else {
				fmt.Println(string(data))
			}
		}
	}()
	return measurements
}
