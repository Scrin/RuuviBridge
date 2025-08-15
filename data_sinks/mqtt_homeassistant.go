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
	DataFormat                string `json:"data_format"`
	Rssi                      *int64 `json:"rssi,omitempty"`
	TxPower                   *int64 `json:"tx_power,omitempty"`
	MeasurementSequenceNumber *int64 `json:"measurement_sequence_number,omitempty"`
	CalibrationInProgress     *bool  `json:"calibration_in_progress,omitempty"`
	ButtonPressedOnBoot       *bool  `json:"button_pressed_on_boot,omitempty"`
	RtcOnBoot                 *bool  `json:"rtc_on_boot,omitempty"`
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
		EntityName:        "temperature",
		UnitOfMeasurement: "°C",
		JsonAttribute:     "temperature",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Humidity != nil,
		DeviceClass:       "humidity",
		EntityName:        "humidity",
		UnitOfMeasurement: "%",
		JsonAttribute:     "humidity",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:            measurement.Pressure != nil,
		DeviceClass:          "pressure",
		EntityName:           "pressure",
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
		EntityName:        "tag battery voltage",
		UnitOfMeasurement: "V",
		JsonAttribute:     "batteryVoltage",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.MovementCounter != nil,
		EntityName:        "movement counter",
		UnitOfMeasurement: "x",
		JsonAttribute:     "movementCounter",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationTotal != nil,
		EntityName:        "total acceleration",
		UnitOfMeasurement: "G",
		JsonAttribute:     "accelerationTotal",
		Icon:              "mdi:axis-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AbsoluteHumidity != nil,
		EntityName:        "absolute humidity",
		UnitOfMeasurement: "g/m³",
		JsonAttribute:     "absoluteHumidity",
		Icon:              "mdi:water",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.DewPoint != nil,
		DeviceClass:       "temperature",
		EntityName:        "dew point",
		UnitOfMeasurement: "°C",
		JsonAttribute:     "dewPoint",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:            measurement.EquilibriumVaporPressure != nil,
		DeviceClass:          "pressure",
		EntityName:           "equilibrium vapor pressure",
		UnitOfMeasurement:    "hPa",
		JsonAttribute:        "equilibriumVaporPressure",
		JsonAttributeMutator: " / 100.0",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AirDensity != nil,
		EntityName:        "air density",
		UnitOfMeasurement: "kg/m³",
		JsonAttribute:     "airDensity",
		Icon:              "mdi:gauge",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromX != nil,
		EntityName:        "acceleration angle from X axis",
		UnitOfMeasurement: "º",
		JsonAttribute:     "accelerationAngleFromX",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromY != nil,
		EntityName:        "acceleration angle from Y axis",
		UnitOfMeasurement: "º",
		JsonAttribute:     "accelerationAngleFromY",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromZ != nil,
		EntityName:        "acceleration angle from Z axis",
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
		EntityName:        "measurement sequence number",
		UnitOfMeasurement: "x",
		JsonAttribute:     "measurementSequenceNumber",
		Icon:              "mdi:counter",
	})
	// New E1 fields
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm10 != nil,
		EntityName:        "PM1.0",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm10",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm25 != nil,
		EntityName:        "PM2.5",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm25",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm40 != nil,
		EntityName:        "PM4.0",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm40",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm100 != nil,
		EntityName:        "PM10.0",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm100",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.CO2 != nil,
		DeviceClass:       "carbon_dioxide",
		EntityName:        "CO₂",
		UnitOfMeasurement: "ppm",
		JsonAttribute:     "co2",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.VOC != nil,
		EntityName:        "VOC index",
		UnitOfMeasurement: "x",
		JsonAttribute:     "voc",
		Icon:              "mdi:chemical-weapon",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.NOX != nil,
		EntityName:        "NOx index",
		UnitOfMeasurement: "x",
		JsonAttribute:     "nox",
		Icon:              "mdi:chemical-weapon",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Illuminance != nil,
		EntityName:        "illuminance",
		UnitOfMeasurement: "lx",
		JsonAttribute:     "illuminance",
		Icon:              "mdi:brightness-5",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.SoundInstant != nil,
		EntityName:        "dBA instant",
		UnitOfMeasurement: "dBA",
		JsonAttribute:     "soundInstant",
		Icon:              "mdi:volume-high",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.SoundAverage != nil,
		EntityName:        "dBA average",
		UnitOfMeasurement: "dBA",
		JsonAttribute:     "soundAverage",
		Icon:              "mdi:volume-medium",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.SoundPeak != nil,
		EntityName:        "dBA peak",
		UnitOfMeasurement: "dBA",
		JsonAttribute:     "soundPeak",
		Icon:              "mdi:volume-high",
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
		Mac:                   measurement.Mac,
		DataFormat:            fmt.Sprintf("%X", measurement.DataFormat),
		CalibrationInProgress: measurement.CalibrationInProgress,
		ButtonPressedOnBoot:   measurement.ButtonPressedOnBoot,
		RtcOnBoot:             measurement.RtcOnBoot,
	})
	if err != nil {
		log.WithError(err).Error("Failed to serialize Home Assistant attribute data")
		return
	}
	client.Publish(confTopicPrefix+"/attributes", 0, false, string(attributesJson))
	client.Publish(confTopicPrefix+"/config", 0, false, string(discoveryJson))
}
