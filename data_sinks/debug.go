package data_sinks

import (
	"encoding/json"

	"github.com/Scrin/RuuviBridge/parser"
	log "github.com/sirupsen/logrus"
)

func Debug() chan<- parser.Measurement {
	log.Info("Starting debug sink")
	measurements := make(chan parser.Measurement, 1024)
	go func() {
		for measurement := range measurements {
			data, err := json.Marshal(measurement)
			if err != nil {
				log.WithError(err).Error("Failed to serialize measurement")
			} else {
				var fields map[string]interface{}
				err = json.Unmarshal(data, &fields)
				if err != nil {
					log.WithError(err).Error("Failed to deserialize measurement")
				} else {
					log.WithFields(fields).Println("Processed measurement")
				}
			}
		}
	}()
	return measurements
}
