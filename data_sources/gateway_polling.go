package data_sources

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type gatewayHistoryTag struct {
	Rssi      int64  `json:"rssi"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
}

// seems to be emitted only if the authentication fails
type gatewayInfo struct {
	GatewayName string `json:"gateway_name"`
}

type gatewayHistory struct {
	Data struct {
		GwMac string                       `json:"gw_mac"`
		Tags  map[string]gatewayHistoryTag `json:"tags"`
	} `json:"data"`
}

func StartGatewayPolling(conf config.GatewayPolling, measurements chan<- parser.Measurement) chan<- bool {
	interval := conf.Interval
	if interval == 0 {
		interval = 10 * time.Second
	}
	logger := log.With().
		Str("target", conf.GatewayUrl).
		Dur("interval", interval).
		Logger()
	logger.Info().Msg("Starting gateway polling")
	stop := make(chan bool)
	go gatewayPoller(conf.GatewayUrl, conf.BearerToken, interval, measurements, stop, logger)
	return stop
}

func gatewayPoller(url string, bearer_token string, interval time.Duration, measurements chan<- parser.Measurement, stop <-chan bool, logger zerolog.Logger) {
	seenTags := make(map[string]int64)
	poll(url, bearer_token, measurements, seenTags, logger)
	for {
		select {
		case <-stop:
			return
		case <-time.After(interval):
			poll(url, bearer_token, measurements, seenTags, logger)
		}
	}
}

func poll(url string, bearer_token string, measurements chan<- parser.Measurement, seenTags map[string]int64, logger zerolog.Logger) {
	req, err := http.NewRequest("GET", url+"/history", nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to construct GET request")
		return
	}

	if bearer_token != "" {
		req.Header.Add("Authorization", "Bearer "+bearer_token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get history from gateway")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read data from gateway")
		return
	}

	var gatewayInfo gatewayInfo
	err = json.Unmarshal(body, &gatewayInfo)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to deserialize gateway data")
		return
	}
	if len(gatewayInfo.GatewayName) > 0 {
		logger.Error().Msg("Failed to authenticate")
		return
	}

	var gatewayHistory gatewayHistory
	err = json.Unmarshal(body, &gatewayHistory)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to deserialize gateway data")
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
