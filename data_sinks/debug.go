package data_sinks

import (
	"encoding/json"

	"github.com/Scrin/RuuviBridge/parser"
	"github.com/rs/zerolog/log"
)

func Debug() chan<- parser.Measurement {
	log.Info().Msg("Starting debug sink")
	measurements := make(chan parser.Measurement, 1024)
	go func() {
		for measurement := range measurements {
			data, err := json.Marshal(measurement)
			if err != nil {
				log.Error().Err(err).Msg("Failed to serialize measurement")
			} else {
				var fields map[string]interface{}
				err = json.Unmarshal(data, &fields)
				if err != nil {
					log.Error().Err(err).Msg("Failed to deserialize measurement")
				} else {
					log.Info().Fields(fields).Msg("Processed measurement")
				}
			}
		}
	}()
	return measurements
}
