package data_sources

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type message struct {
	GwMac  string        `json:"gw_mac"`
	Rssi   int64         `json:"rssi"`
	Aoa    []interface{} `json:"aoa"`
	Gwts   string        `json:"gwts"`
	Ts     string        `json:"ts"`
	Data   string        `json:"data"`
	Coords string        `json:"coords"`
}

func StartMQTTListener(conf config.MQTTListener, measurements chan<- parser.Measurement) chan<- bool {
	address := conf.BrokerAddress
	if address == "" {
		address = "localhost"
	}
	port := conf.BrokerPort
	if port == 0 {
		port = 1883
	}
	server := fmt.Sprintf("tcp://%s:%d", address, port)

	log.WithFields(log.Fields{"target": server}).Info("Starting MQTT subscriber")

	messagePubHandler := func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		var message message
		err := json.Unmarshal(msg.Payload(), &message)
		if err != nil {
			log.Error(err)
			return
		}

		mac := topic[strings.LastIndex(topic, "/")+1:]
		timestamp, _ := strconv.ParseInt(message.Ts, 10, 64)

		measurement, ok := parser.Parse(message.Data)
		if ok {
			measurement.Mac = mac
			measurement.Rssi = &message.Rssi
			measurement.Timestamp = &timestamp
			measurements <- measurement
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetClientID(conf.ClientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := client.Subscribe(conf.TopicPrefix+"/+", 0, messagePubHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	stop := make(chan bool)
	go func() {
		<-stop
		client.Disconnect(0)
	}()
	return stop
}
