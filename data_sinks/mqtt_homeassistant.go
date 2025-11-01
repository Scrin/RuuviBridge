package data_sinks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
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
	AvailabilityTopic   string                       `json:"availability_topic,omitempty"`
	PayloadAvailable    string                       `json:"payload_available,omitempty"`
	PayloadNotAvailable string                       `json:"payload_not_available,omitempty"`
	EntityCategory      string                       `json:"entity_category,omitempty"`
	Device              homeassistantDiscoveryDevice `json:"device"`
}

type homeassistantDiscoveryAttributes struct {
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
	StateClass           string
	EntityCategory       string
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
		EntityName:        "Acceleration X",
		UnitOfMeasurement: "g",
		JsonAttribute:     "accelerationX",
		Icon:              "mdi:axis-x-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationY != nil,
		EntityName:        "Acceleration Y",
		UnitOfMeasurement: "g",
		JsonAttribute:     "accelerationY",
		Icon:              "mdi:axis-y-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationZ != nil,
		EntityName:        "Acceleration Z",
		UnitOfMeasurement: "g",
		JsonAttribute:     "accelerationZ",
		Icon:              "mdi:axis-z-arrow",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.BatteryVoltage != nil,
		DeviceClass:       "voltage",
		EntityName:        "Battery voltage",
		UnitOfMeasurement: "V",
		JsonAttribute:     "batteryVoltage",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.MovementCounter != nil,
		EntityName:        "Movement counter",
		UnitOfMeasurement: "x",
		JsonAttribute:     "movementCounter",
		StateClass:        "total_increasing",
		EntityCategory:    "diagnostic",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationTotal != nil,
		EntityName:        "Total acceleration",
		UnitOfMeasurement: "g",
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
		UnitOfMeasurement: "°",
		JsonAttribute:     "accelerationAngleFromX",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromY != nil,
		EntityName:        "Acceleration angle from Y axis",
		UnitOfMeasurement: "°",
		JsonAttribute:     "accelerationAngleFromY",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.AccelerationAngleFromZ != nil,
		EntityName:        "Acceleration angle from Z axis",
		UnitOfMeasurement: "°",
		JsonAttribute:     "accelerationAngleFromZ",
		Icon:              "mdi:angle-acute",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Rssi != nil,
		DeviceClass:       "signal_strength",
		EntityName:        "RSSI",
		UnitOfMeasurement: "dBm",
		JsonAttribute:     "rssi",
		Icon:              "mdi:signal-variant",
		EntityCategory:    "diagnostic",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.TxPower != nil,
		EntityName:        "TX power",
		UnitOfMeasurement: "dBm",
		JsonAttribute:     "txPower",
		Icon:              "mdi:signal-variant",
		EntityCategory:    "diagnostic",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.MeasurementSequenceNumber != nil,
		EntityName:        "Measurement sequence number",
		UnitOfMeasurement: "x",
		JsonAttribute:     "measurementSequenceNumber",
		Icon:              "mdi:counter",
		StateClass:        "total_increasing",
		EntityCategory:    "diagnostic",
	})
	// New E1 fields
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm1p0 != nil,
		DeviceClass:       "pm1",
		EntityName:        "PM1.0",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm1p0",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm2p5 != nil,
		DeviceClass:       "pm25",
		EntityName:        "PM2.5",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm2p5",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm4p0 != nil,
		EntityName:        "PM4.0",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm4p0",
		Icon:              "mdi:molecule",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Pm10p0 != nil,
		DeviceClass:       "pm10",
		EntityName:        "PM10",
		UnitOfMeasurement: "µg/m³",
		JsonAttribute:     "pm10p0",
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
		Icon:              "mdi:molecule",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.NOX != nil,
		EntityName:        "NOx index",
		UnitOfMeasurement: "x",
		JsonAttribute:     "nox",
		Icon:              "mdi:molecule",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.Illuminance != nil,
		DeviceClass:       "illuminance",
		EntityName:        "Illuminance",
		UnitOfMeasurement: "lx",
		JsonAttribute:     "illuminance",
		Icon:              "mdi:brightness-5",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.SoundInstant != nil,
		DeviceClass:       "sound_pressure",
		EntityName:        "Sound level (instant, A-weighted)",
		UnitOfMeasurement: "dB",
		JsonAttribute:     "soundInstant",
		Icon:              "mdi:volume-medium",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.SoundAverage != nil,
		DeviceClass:       "sound_pressure",
		EntityName:        "Sound level (average, A-weighted)",
		UnitOfMeasurement: "dB",
		JsonAttribute:     "soundAverage",
		Icon:              "mdi:volume-medium",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:         measurement.SoundPeak != nil,
		DeviceClass:       "sound_pressure",
		EntityName:        "Sound level (peak, A-weighted)",
		UnitOfMeasurement: "dB",
		JsonAttribute:     "soundPeak",
		Icon:              "mdi:volume-high",
	})
	publishHomeAssistantDiscovery(client, conf, measurement, homeassistantDiscoveryConfig{
		Available:     measurement.AirQualityIndex != nil,
		DeviceClass:   "aqi",
		EntityName:    "Air quality index",
		JsonAttribute: "airQualityIndex",
	})
}

func publishHomeAssistantDiscovery(client mqtt.Client, conf config.MQTTPublisher, measurement parser.Measurement, disco homeassistantDiscoveryConfig) {
	id := fmt.Sprintf("ruuvitag_%s_%s", strings.ReplaceAll(measurement.Mac, ":", ""), disco.JsonAttribute)
	confTopicPrefix := fmt.Sprintf("%s/sensor/%s", conf.HomeassistantDiscoveryPrefix, id)
	if !disco.Available {
		client.Publish(confTopicPrefix+"/config", 0, conf.RetainMessages, "")
		client.Publish(confTopicPrefix+"/attributes", 0, conf.RetainMessages, "")
		return
	}
	var name string
	if measurement.Name != nil {
		name = *measurement.Name
	} else {
		name = fmt.Sprintf("RuuviTag %s", measurement.Mac)
	}
	stateClass := disco.StateClass
	if stateClass == "" {
		stateClass = "measurement"
	}
	discoveryJson, err := json.Marshal(homeassistantDiscovery{
		UniqueID:            id,
		DeviceClass:         disco.DeviceClass,
		StateTopic:          conf.TopicPrefix + "/" + measurement.Mac,
		StateClass:          stateClass,
		JsonAttributesTopic: confTopicPrefix + "/attributes",
		Name:                disco.EntityName,
		UnitOfMeasurement:   disco.UnitOfMeasurement,
		ValueTemplate:       fmt.Sprintf("{{ (value_json.%s%s) | round(2) }}", disco.JsonAttribute, disco.JsonAttributeMutator),
		Icon:                disco.Icon,
		AvailabilityTopic:   conf.LWTTopic,
		PayloadAvailable:    conf.LWTOnlinePayload,
		PayloadNotAvailable: conf.LWTOfflinePayload,
		EntityCategory:      disco.EntityCategory,
		Device: homeassistantDiscoveryDevice{
			Identifiers:  []string{measurement.Mac},
			Name:         name,
			Model:        "RuuviTag",
			Manufacturer: "Ruuvi",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize Home Assistant discovery data")
		return
	}
	attributesJson, err := json.Marshal(homeassistantDiscoveryAttributes{
		Mac:                   measurement.Mac,
		DataFormat:            fmt.Sprintf("%X", measurement.DataFormat),
		CalibrationInProgress: measurement.CalibrationInProgress,
		ButtonPressedOnBoot:   measurement.ButtonPressedOnBoot,
		RtcOnBoot:             measurement.RtcOnBoot,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize Home Assistant attribute data")
		return
	}
	client.Publish(confTopicPrefix+"/attributes", 0, conf.RetainMessages, string(attributesJson))
	client.Publish(confTopicPrefix+"/config", 0, conf.RetainMessages, string(discoveryJson))
}
