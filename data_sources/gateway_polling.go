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

// seems to be emitted only if the authentication fails
type gatewayInfo struct {
	Success     bool   `json:"success"`
	GatewayName string `json:"gateway_name"`
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
	log := log.WithFields(log.Fields{
		"target":   conf.GatewayUrl,
		"interval": interval,
	})
	log.Info("Starting gateway polling")
	stop := make(chan bool)
	go gatewayPoller(conf.GatewayUrl, conf.BearerToken, interval, measurements, stop, log)
	return stop
}

func gatewayPoller(url string, bearer_token string, interval time.Duration, measurements chan<- parser.Measurement, stop <-chan bool, log *log.Entry) {
	seenTags := make(map[string]int64)
	poll(url, bearer_token, measurements, seenTags, log)
	for {
		select {
		case <-stop:
			return
		case <-time.After(interval):
			poll(url, bearer_token, measurements, seenTags, log)
		}
	}
}

func poll(url string, bearer_token string, measurements chan<- parser.Measurement, seenTags map[string]int64, log *log.Entry) {
	req, err := http.NewRequest("GET", url+"/history", nil)
	if err != nil {
		log.WithError(err).Error("Failed to construct GET request")
		return
	}

	if bearer_token != "" {
		req.Header.Add("Authorization", "Bearer "+bearer_token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to get history from gateway")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read data from gateway")
		return
	}

	// initialize Success so we can detect if the field was populated
	gatewayInfo := gatewayInfo{Success: true, GatewayName: ""}
	err = json.Unmarshal(body, &gatewayInfo)
	if err != nil {
		log.WithError(err).Error("Failed to deserialize gateway data")
		return
	}
	if !gatewayInfo.Success {
		log.Error("Failed to authenticate")
		return
	}

	var gatewayHistory gatewayHistory
	err = json.Unmarshal(body, &gatewayHistory)
	if err != nil {
		log.WithError(err).Error("Failed to deserialize gateway data")
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
