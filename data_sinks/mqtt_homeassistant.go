package data_sinks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type homeassistantDiscoveryDevice struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Model        string   `json:"model"`
	Manufacturer string   `json:"manufacturer"`
}

type homeassistantDiscovery struct {
	UniqueID            string                       `json:"unique_id"`
	DeviceClass         string                       `json:"device_class,omitempty"`
	StateTopic          string                       `json:"state_topic"`
	StateClass          string                       `json:"state_class"`
	JsonAttributesTopic string                       `json:"json_attributes_topic"`
	Name                string                       `json:"name,omitempty"`
	UnitOfMeasurement   string                       `json:"unit_of_measurement"`
	ValueTemplate       string                       `json:"value_template"`
	Icon                string                       `json:"icon,omitempty"`
	Device              homeassistantDiscoveryDevice `json:"device"`
}

type homeassistantAttributes struct {
	Mac                       string `json:"mac"`
	DataFormat                int64  `json:"data_format"`
	Rssi                      *int64 `json:"rssi,omitempty"`
	TxPower                   *int64 `json:"tx_power,omitempty"`
	MeasurementSequenceNumber *int64 `json:"measurement_sequence_number,omitempty"`
}

type homeassistantDiscoveryConfig struct {
	Available            bool
	DeviceClass          string
	EntityName           string
	UnitOfMeasurement    string
	JsonAttribute        string
	JsonAttributeMutator string
	Icon                 string
}

func publishHomeAssistantDiscoveries(client mqtt.Client, conf config.MQTTPublisher, measurement parser.Measurement) {
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Temperature != nil,
		DeviceClass:       "temperature",
		EntityName:        "Temperature",
		UnitOfMeasurement: "°C",
		JsonAttribute:     "temperature",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Humidity != nil,
		DeviceClass:       "humidity",
		EntityName:        "Humidity",
		UnitOfMeasurement: "%",
		JsonAttribute:     "humidity",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:            measurement.Pressure != nil,
		DeviceClass:          "pressure",
		EntityName:           "Pressure",
		UnitOfMeasurement:    "hPa",
		JsonAttribute:        "pressure",
		JsonAttributeMutator: " / 100.0",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationX != nil,
		EntityName:        "X acceleration",
		UnitOfMeasurement: "G",
		JsonAttribute:     "accelerationX",
		Icon:              "mdi:axis-x-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationY != nil,
		EntityName:        "Y acceleration",
		UnitOfMeasurement: "G",
		JsonAttribute:     "accelerationY",
		Icon:              "mdi:axis-y-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationZ != nil,
		EntityName:        "Z acceleration",
		UnitOfMeasurement: "G",
		JsonAttribute:     "accelerationZ",
		Icon:              "mdi:axis-z-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.BatteryVoltage != nil,
		DeviceClass:       "voltage",
		EntityName:        "Tag battery voltage",
		UnitOfMeasurement: "V",
		JsonAttribute:     "batteryVoltage",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.MovementCounter != nil,
		EntityName:        "Movement counter",
		UnitOfMeasurement: "x",
		JsonAttribute:     "movementCounter",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationTotal != nil,
		EntityName:        "Total acceleration",
		UnitOfMeasurement: "G",
		JsonAttribute:     "accelerationTotal",
		Icon:              "mdi:axis-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AbsoluteHumidity != nil,
		EntityName:        "Absolute humidity",
		UnitOfMeasurement: "g/m³",
		JsonAttribute:     "absoluteHumidity",
		Icon:              "mdi:water",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.DewPoint != nil,
		DeviceClass:       "temperature",
		EntityName:        "Dew point",
		UnitOfMeasurement: "°C",
		JsonAttribute:     "dewPoint",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:            measurement.EquilibriumVaporPressure != nil,
		DeviceClass:          "pressure",
		EntityName:           "Equilibrium vapor pressure",
		UnitOfMeasurement:    "hPa",
		JsonAttribute:        "equilibriumVaporPressure",
		JsonAttributeMutator: " / 100.0",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AirDensity != nil,
		EntityName:        "Air density",
		UnitOfMeasurement: "kg/m³",
		JsonAttribute:     "airDensity",
		Icon:              "mdi:gauge",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromX != nil,
		EntityName:        "Acceleration angle from X axis",
		UnitOfMeasurement: "º",
		JsonAttribute:     "accelerationAngleFromX",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromY != nil,
		EntityName:        "Acceleration angle from Y axis",
		UnitOfMeasurement: "º",
		JsonAttribute:     "accelerationAngleFromY",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromZ != nil,
		EntityName:        "Acceleration angle from Z axis",
		UnitOfMeasurement: "º",
		JsonAttribute:     "accelerationAngleFromZ",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Rssi != nil,
		EntityName:        "RSSI",
		UnitOfMeasurement: "dBm",
		JsonAttribute:     "rssi",
		Icon:              "mdi:signal-variant",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.TxPower != nil,
		EntityName:        "TX power",
		UnitOfMeasurement: "dBm",
		JsonAttribute:     "txPower",
		Icon:              "mdi:signal-variant",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.MeasurementSequenceNumber != nil,
		EntityName:        "Measurement sequence number",
		UnitOfMeasurement: "x",
		JsonAttribute:     "measurementSequenceNumber",
		Icon:              "mdi:counter",
	})
}

func publishHomeAssistantDiscovery(client mqtt.Client, conf config.MQTTPublisher, measurement parser.Measurement, disco homeassistantDiscoveryConfig) {
	id := fmt.Sprintf("ruuvitag_%s_%s", strings.ReplaceAll(measurement.Mac, ":", ""), disco.JsonAttribute)
	confTopicPrefix := fmt.Sprintf("%s/sensor/%s", conf.HomeassistantDiscoveryPrefix, id)
	if !disco.Available {
		client.Publish(confTopicPrefix+"/config", 0, false, "")
		client.Publish(confTopicPrefix+"/attributes", 0, false, "")
		return
	}
	var name string
	if measurement.Name != nil {
		name = *measurement.Name
	} else {
		name = fmt.Sprintf("RuuviTag %s", measurement.Mac)
	}
	discoveryJson, err := json.Marshal(homeassistantDiscovery{
		UniqueID:            id,
		DeviceClass:         disco.DeviceClass,
		StateTopic:          conf.TopicPrefix + "/" + measurement.Mac,
		StateClass:          "measurement",
		JsonAttributesTopic: confTopicPrefix + "/attributes",
		Name:                disco.EntityName,
		UnitOfMeasurement:   disco.UnitOfMeasurement,
		ValueTemplate:       fmt.Sprintf("{{ (value_json.%s%s) | round(2) }}", disco.JsonAttribute, disco.JsonAttributeMutator),
		Icon:                disco.Icon,
		Device: homeassistantDiscoveryDevice{
			Identifiers:  []string{measurement.Mac},
			Name:         name,
			Model:        "RuuviTag",
			Manufacturer: "Ruuvi",
		},
	})
	if err != nil {
		log.WithError(err).Error("Failed to serialize Home Assistant discovery data")
		return
	}
	attributesJson, err := json.Marshal(homeassistantAttributes{
		Mac:        measurement.Mac,
		DataFormat: measurement.DataFormat,
	})
	if err != nil {
		log.WithError(err).Error("Failed to serialize Home Assistant attribute data")
		return
	}
	client.Publish(confTopicPrefix+"/attributes", 0, false, string(attributesJson))
	client.Publish(confTopicPrefix+"/config", 0, false, string(discoveryJson))
}
