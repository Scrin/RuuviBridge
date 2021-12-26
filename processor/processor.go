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
	if config.GatewayPolling != nil {
		stop := data_sources.StartGatewayPolling(*config.GatewayPolling, measurements)
		defer func() { stop <- true }()
	}
	if config.MQTTListener != nil {
		stop := data_sources.StartMQTTListener(*config.MQTTListener, measurements)
		defer func() { stop <- true }()
	}

	fmt.Println("Starting data sinks...")
	if config.Debug {
		sinks = append(sinks, data_sinks.Debug())
	}
	if config.InfluxDBPublisher != nil {
		sinks = append(sinks, data_sinks.InfluxDB(*config.InfluxDBPublisher))
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
