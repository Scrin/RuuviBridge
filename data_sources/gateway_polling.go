package data_sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
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

func StartGatewayPolling(polling config.GatewayPolling, measurements chan<- parser.Measurement) chan<- bool {
	interval := polling.Interval
	if interval == 0 {
		interval = 10 * time.Second
	}
	fmt.Printf("Starting gateway polling at %s every %s\n", polling.GatewayUrl, polling.Interval)
	stop := make(chan bool)
	go gatewayPoller(polling.GatewayUrl, interval, measurements, stop)
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
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("Got %d status code, expected code 200", resp.StatusCode)
		return
	}
	var gatewayHistory gatewayHistory
	err = json.Unmarshal(body, &gatewayHistory)
	if err != nil {
		fmt.Println(err)
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
