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
	DeviceClass         string `json:"device_class"`
	StateTopic          string `json:"state_topic"`
	JsonAttributesTopic string `json:"json_attributes_topic"`
	Name                string `json:"name,omitempty"`
	UnitOfMeasurement   string `json:"unit_of_measurement"`
	ValueTemplate       string `json:"value_template"`
}

type homeassistantAttributes struct {
	Mac        string `json:"mac"`
	DataFormat int64  `json:"data_format"`
}

func publishHomeAssistantDiscoveries(client mqtt.Client, conf config.MQTTPublisher, measurement parser.Measurement) {
	name := "RuuviTag " + measurement.Mac
	if measurement.Name != nil {
		name = *measurement.Name
	}
	attrs := homeassistantAttributes{
		Mac:        measurement.Mac,
		DataFormat: measurement.DataFormat,
	}
	topicPrefix := fmt.Sprintf("%s/sensor/ruuvitag_%s_temperature", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", ""))
	publishHomeAssistantDiscovery(client,
		topicPrefix,
		homeassistantDiscovery{
			DeviceClass:         "temperature",
			StateTopic:          conf.TopicPrefix + "/" + measurement.Mac,
			JsonAttributesTopic: topicPrefix + "/attributes",
			Name:                name + " temperature",
			UnitOfMeasurement:   "ºC",
			ValueTemplate:       "{{ value_json.temperature }}",
		}, attrs)
	topicPrefix = fmt.Sprintf("%s/sensor/ruuvitag_%s_humidity", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", ""))
	publishHomeAssistantDiscovery(client,
		topicPrefix,
		homeassistantDiscovery{
			DeviceClass:         "humidity",
			StateTopic:          conf.TopicPrefix + "/" + measurement.Mac,
			JsonAttributesTopic: topicPrefix + "/attributes",
			Name:                name + " humidity",
			UnitOfMeasurement:   "%",
			ValueTemplate:       "{{ value_json.humidity }}",
		}, attrs)
	topicPrefix = fmt.Sprintf("%s/sensor/ruuvitag_%s_pressure", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", ""))
	publishHomeAssistantDiscovery(client,
		topicPrefix,
		homeassistantDiscovery{
			DeviceClass:         "pressure",
			StateTopic:          conf.TopicPrefix + "/" + measurement.Mac,
			JsonAttributesTopic: topicPrefix + "/attributes",
			Name:                name + " pressure",
			UnitOfMeasurement:   "hPa",
			ValueTemplate:       "{{ value_json.pressure / 100.0 }}",
		}, attrs)
	topicPrefix = fmt.Sprintf("%s/sensor/ruuvitag_%s_battery", conf.HomeassistantDiscoveryPrefix, strings.ReplaceAll(measurement.Mac, ":", ""))
	publishHomeAssistantDiscovery(client,
		topicPrefix,
		homeassistantDiscovery{
			DeviceClass:         "voltage",
			StateTopic:          conf.TopicPrefix + "/" + measurement.Mac,
			JsonAttributesTopic: topicPrefix + "/attributes",
			Name:                name + " battery voltage",
			UnitOfMeasurement:   "V",
			ValueTemplate:       "{{ value_json.batteryVoltage }}",
		}, attrs)
}

func publishHomeAssistantDiscovery(client mqtt.Client, topicPrefix string, disco homeassistantDiscovery, attrs homeassistantAttributes) {
	discoveryJson, err := json.Marshal(disco)
	if err != nil {
		fmt.Println(err)
		return
	}
	attributesJson, err := json.Marshal(attrs)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Publish(topicPrefix+"/config", 0, false, string(discoveryJson))
	client.Publish(topicPrefix+"/attributes", 0, false, string(attributesJson))
}
