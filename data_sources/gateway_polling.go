package data_sources

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	log "github.com/sirupsen/logrus"
)

type gatewayHistoryTag struct {
	Rssi      int64  `json:"rssi"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

type gatewayHistory struct {
	Data struct {
		Coordinates string                       `json:"coordinates"`
		Timestamp   string                       `json:"timestamp"`
		GwMac       string                       `json:"gw_mac"`
		Tags        map[string]gatewayHistoryTag `json:"tags"`
	} `json:"data"`
}

func StartGatewayPolling(conf config.GatewayPolling, measurements chan<- parser.Measurement) chan<- bool {
	interval := conf.Interval
	if interval == 0 {
		interval = 10 * time.Second
	}
	log.WithFields(log.Fields{"target": conf.GatewayUrl, "interval": conf.Interval}).Info("Starting gateway polling")
	stop := make(chan bool)
	go gatewayPoller(conf.GatewayUrl, interval, measurements, stop)
	return stop
}

func gatewayPoller(url string, interval time.Duration, measurements chan<- parser.Measurement, stop <-chan bool) {
	seenTags := make(map[string]int64)
	poll(url, measurements, seenTags)
	for {
		select {
		case <-stop:
			return
		case <-time.After(interval):
			poll(url, measurements, seenTags)
		}
	}
}

func poll(url string, measurements chan<- parser.Measurement, seenTags map[string]int64) {
	resp, err := http.Get(url + "/history")
	if err != nil {
		log.Error("Failed to get history from gateway: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read data from gateway: ", err)
		return
	}
	var gatewayHistory gatewayHistory
	err = json.Unmarshal(body, &gatewayHistory)
	if err != nil {
		log.Error(err)
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
