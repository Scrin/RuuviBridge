package data_sinks

import (
	"encoding/json"
	"fmt"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MQTT(conf config.MQTTPublisher) chan<- parser.Measurement {
	address := conf.BrokerAddress
	if address == "" {
		address = "localhost"
	}
	port := conf.BrokerPort
	if port == 0 {
		port = 1883
	}
	server := fmt.Sprintf("tcp://%s:%d", address, port)
	fmt.Println("Starting MQTT sink")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetClientID(conf.ClientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	measurements := make(chan parser.Measurement)
	go func() {
		for measurement := range measurements {
			data, err := json.Marshal(measurement)
			if err != nil {
				fmt.Println(err)
			} else {
				client.Publish(conf.TopicPrefix+"/"+measurement.Mac, 0, false, string(data))
				if conf.HomeassistantDiscoveryPrefix != "" {
					publishHomeAssistantDiscoveries(client, conf, measurement)
				}
			}
		}
	}()
	return measurements
}
