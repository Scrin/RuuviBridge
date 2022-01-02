package data_sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	log "github.com/sirupsen/logrus"
)

func StartHTTPListener(conf config.HTTPListener, measurements chan<- parser.Measurement) chan<- bool {
	port := conf.Port
	if port == 0 {
		port = 8080
	}
	log.WithFields(log.Fields{"port": port}).Info("Starting http listener")

	seenTags := make(map[string]int64)

	handlerFunc := func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"path": req.URL.Path,
			}).WithError(err).Error("Failed to read request body")
			return
		}
		req.Body.Close()
		log := log.WithFields(log.Fields{
			"path": req.URL.Path,
			"body": string(body),
		})
		log.Trace("Received a http call")

		var gatewayHistory gatewayHistory
		err = json.Unmarshal(body, &gatewayHistory)
		if err != nil {
			log.WithError(err).Error("Failed to deserialize http listener data")
			return
		}

		for mac, data := range gatewayHistory.Data.Tags {
			mac = strings.ToUpper(mac)
			timestamp, _ := strconv.ParseInt(data.Timestamp, 10, 64)
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
