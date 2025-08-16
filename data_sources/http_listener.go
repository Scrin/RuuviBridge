package data_sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	"github.com/rs/zerolog/log"
)

func StartHTTPListener(conf config.HTTPListener, measurements chan<- parser.Measurement) chan<- bool {
	port := conf.Port
	if port == 0 {
		port = 8080
	}
	log.Info().Int("port", port).Msg("Starting http listener")

	seenTags := make(map[string]int64)

	handlerFunc := func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error().Str("path", req.URL.Path).Err(err).Msg("Failed to read request body")
			return
		}
		req.Body.Close()
		logger := log.With().
			Str("path", req.URL.Path).
			Str("body", string(body)).
			Logger()
		logger.Trace().Msg("Received a http call")

		var gatewayHistory gatewayHistory
		err = json.Unmarshal(body, &gatewayHistory)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to deserialize http listener data")
			return
		}

		for mac, data := range gatewayHistory.Data.Tags {
			mac = strings.ToUpper(mac)
			timestamp := data.Timestamp
			if seenTags[mac] == timestamp {
				continue
			}
			seenTags[mac] = timestamp
			measurement, ok := parser.Parse(data.Data)
			if ok {
				measurement.Mac = mac
				measurement.Rssi = &data.Rssi
				measurement.Timestamp = &timestamp
				measurements <- measurement
			}
		}
	}

	serverMuxA := http.NewServeMux()
	serverMuxA.HandleFunc("/", handlerFunc)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), serverMuxA)

	stop := make(chan bool)
	go func() {
		<-stop
	}()
	return stop
}
