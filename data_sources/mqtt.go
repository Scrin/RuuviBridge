package data_sources

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

type message struct {
	GwMac  string        `json:"gw_mac"`
	Rssi   int64         `json:"rssi"`
	Aoa    []interface{} `json:"aoa"`
	Gwts   interface{}   `json:"gwts"`
	Ts     interface{}   `json:"ts"`
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
	server := conf.BrokerUrl
	if server == "" {
		server = fmt.Sprintf("tcp://%s:%d", address, port)
	}
	subscription := conf.TopicPrefix + "/+"
	log := log.With().
		Str("target", server).
		Str("topic_prefix", conf.TopicPrefix).
		Str("mqtt_subscription", subscription).
		Logger()

	log.Info().Msg("Starting MQTT subscriber")

	messagePubHandler := func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		var message message
		err := json.Unmarshal(msg.Payload(), &message)
		if err != nil {
			log.Error().Err(err).Msg("Failed to deserialize MQTT message")
			return
		}

		mac := strings.ToUpper(topic[strings.LastIndex(topic, "/")+1:])
		timestamp, _ := strconv.ParseInt(fmt.Sprint(message.Ts), 10, 64)

		measurement, ok := parser.Parse(message.Data)
		if ok {
			measurement.Mac = mac
			measurement.Rssi = &message.Rssi
			measurement.Timestamp = &timestamp
			measurements <- measurement
		}
	}
	clientID := conf.ClientID
	if clientID == "" {
		clientID = "RuuviBridgeListener"
	}
	opts := mqtt.NewClientOptions()
	opts.SetCleanSession(false)
	opts.AddBroker(server)
	opts.SetClientID(clientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	if conf.LWTTopic != "" {
		payload := conf.LWTOfflinePayload
		if payload == "" {
			payload = "{\"state\":\"offline\"}"
		}
		opts.SetWill(conf.LWTTopic, payload, 0, true)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Failed to connect to MQTT")
	}
	if token := client.Subscribe(subscription, 0, messagePubHandler); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Failed to subscribe to MQTT topic")
	}
	if conf.LWTTopic != "" {
		payload := conf.LWTOnlinePayload
		if payload == "" {
			payload = "{\"state\":\"online\"}"
		}
		client.Publish(conf.LWTTopic, 0, true, payload)
	}
	stop := make(chan bool)
	go func() {
		<-stop
		client.Disconnect(0)
	}()
	return stop
}
