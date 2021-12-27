package data_sinks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type homeassistantDiscovery struct {
	DeviceClass       string `json:"device_class"`
	StateTopic        string `json:"state_topic"`
	Name              string `json:"name,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	ValueTemplate     string `json:"value_template"`
}

func publishHomeAssistantDiscovery(client mqtt.Client, topic string, disco homeassistantDiscovery) {
	data, err := json.Marshal(disco)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Publish(topic, 0, false, string(data))
}

func publishHomeAssistantDiscoveries(client mqtt.Client, conf config.MQTTPublisher, measurement parser.Measurement) {
	name := "RuuviTag " + measurement.Mac
	if measurement.Name != nil {
		name = *measurement.Name
	}
	publishHomeAssistantDiscovery(client,
		fmt.Sprintf("%s/sensor/ruuvitag_%s_temperature/config", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", "")),
		homeassistantDiscovery{
			DeviceClass:       "temperature",
			StateTopic:        conf.TopicPrefix + "/" + measurement.Mac,
			Name:              name + " temperature",
			UnitOfMeasurement: "ºC",
			ValueTemplate:     "{{ value_json.temperature }}",
		})
	publishHomeAssistantDiscovery(client,
		fmt.Sprintf("%s/sensor/ruuvitag_%s_humidity/config", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", "")),
		homeassistantDiscovery{
			DeviceClass:       "humidity",
			StateTopic:        conf.TopicPrefix + "/" + measurement.Mac,
			Name:              name + " humidity",
			UnitOfMeasurement: "%",
			ValueTemplate:     "{{ value_json.humidity }}",
		})
	publishHomeAssistantDiscovery(client,
		fmt.Sprintf("%s/sensor/ruuvitag_%s_pressure/config", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", "")),
		homeassistantDiscovery{
			DeviceClass:       "pressure",
			StateTopic:        conf.TopicPrefix + "/" + measurement.Mac,
			Name:              name + " pressure",
			UnitOfMeasurement: "hPa",
			ValueTemplate:     "{{ value_json.pressure / 100.0 }}",
		})
}

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
