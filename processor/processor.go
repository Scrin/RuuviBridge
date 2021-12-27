package processor

import (
	"fmt"
	"strings"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/data_sinks"
	"github.com/Scrin/RuuviBridge/data_sources"
	"github.com/Scrin/RuuviBridge/parser"
	"github.com/Scrin/RuuviBridge/value_calculator"
)

func Run(config config.Config) {
	measurements := make(chan parser.Measurement)
	var sinks []chan<- parser.Measurement
	extendedValues := true // default
	if config.Processing != nil {
		processing := config.Processing
		if processing.ExtendedValues != nil {
			extendedValues = *processing.ExtendedValues
		}
	}

	fmt.Println("Starting data sources...")
	datasourcesStarted := false
	if config.GatewayPolling != nil && (config.GatewayPolling.Enabled == nil || *config.GatewayPolling.Enabled) {
		stop := data_sources.StartGatewayPolling(*config.GatewayPolling, measurements)
		defer func() { stop <- true }()
		datasourcesStarted = true
	}
	if config.MQTTListener != nil && (config.MQTTListener.Enabled == nil || *config.MQTTListener.Enabled) {
		stop := data_sources.StartMQTTListener(*config.MQTTListener, measurements)
		defer func() { stop <- true }()
		datasourcesStarted = true
	}
	if !datasourcesStarted {
		fmt.Println("No datasources configured! Please check the config.")
	}

	fmt.Println("Starting data sinks...")
	datasinksStarted := false
	if config.Debug {
		sinks = append(sinks, data_sinks.Debug())
		datasinksStarted = true
	}
	if config.InfluxDBPublisher != nil && (config.InfluxDBPublisher.Enabled == nil || *config.InfluxDBPublisher.Enabled) {
		sinks = append(sinks, data_sinks.InfluxDB(*config.InfluxDBPublisher))
		datasinksStarted = true
	}
	if config.Prometheus != nil && (config.Prometheus.Enabled == nil || *config.Prometheus.Enabled) {
		sinks = append(sinks, data_sinks.Prometheus(*config.Prometheus))
		datasinksStarted = true
	}
	if config.MQTTPublisher != nil && (config.MQTTPublisher.Enabled == nil || *config.MQTTPublisher.Enabled) {
		sinks = append(sinks, data_sinks.MQTT(*config.MQTTPublisher))
		datasinksStarted = true
	}
	if !datasinksStarted {
		fmt.Println("No data consumers/sinks configured! Please check the config.")
	}

	fmt.Println("Starting processing...")
	for measurement := range measurements {
		if extendedValues {
			value_calculator.CalcExtendedValues(&measurement)
		}
		name := config.TagNames[strings.ReplaceAll(measurement.Mac, ":", "")]
		if name != "" {
			measurement.Name = &name
		}
		for _, sink := range sinks {
			sink <- measurement
		}
	}
}
