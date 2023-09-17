package processor

import (
	"strings"

	"github.com/Scrin/RuuviBridge/common/version"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/data_sinks"
	"github.com/Scrin/RuuviBridge/data_sources"
	"github.com/Scrin/RuuviBridge/parser"
	"github.com/Scrin/RuuviBridge/value_calculator"
	log "github.com/sirupsen/logrus"
)

func Run(config config.Config) {
	log.WithFields(log.Fields{"version": version.Version}).Info("RuuviBridge starting up")
	measurements := make(chan parser.Measurement, 1024)
	var sinks []chan<- parser.Measurement

	extendedValues := true // default
	filterMap := make(map[string]interface{})
	allowlist := false
	denylist := false
	namedOnly := false
	if config.Processing != nil {
		processing := config.Processing
		if processing.ExtendedValues != nil {
			extendedValues = *processing.ExtendedValues
		}
		switch processing.FilterMode {
		case "allowlist":
			allowlist = true
			if len(config.Processing.FilterList) == 0 {
				log.Fatal("filter_mode configured as allowlist but no allowed tags configured!")
			}
		case "denylist":
			denylist = true
			if len(config.Processing.FilterList) == 0 {
				log.Fatal("filter_mode configured as denylist but no denied tags configured!")
			}
		case "named":
			namedOnly = true
			if len(config.TagNames) == 0 {
				log.Fatal("filter_mode configured as named but no tag names configured!")
			}
		case "none":
		default:
			log.Fatal("Unrecognized filter_mode: ", processing.FilterMode)
		}
		for _, mac := range config.Processing.FilterList {
			formattedMac := strings.ToUpper(strings.ReplaceAll(mac, ":", ""))
			filterMap[formattedMac] = struct{}{}
		}
	}

	log.Info("Starting data sources")
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
	if config.HTTPListener != nil && (config.HTTPListener.Enabled == nil || *config.HTTPListener.Enabled) {
		stop := data_sources.StartHTTPListener(*config.HTTPListener, measurements)
		defer func() { stop <- true }()
		datasourcesStarted = true
	}
	if !datasourcesStarted {
		log.Fatal("No datasources configured! Please check the config.")
	}

	log.Info("Starting data sinks")
	datasinksStarted := false
	if config.Debug {
		sinks = append(sinks, data_sinks.Debug())
		datasinksStarted = true
	}
	if config.InfluxDBPublisher != nil && (config.InfluxDBPublisher.Enabled == nil || *config.InfluxDBPublisher.Enabled) {
		sinks = append(sinks, data_sinks.InfluxDB(*config.InfluxDBPublisher))
		datasinksStarted = true
	}
	if config.InfluxDB3Publisher != nil && (config.InfluxDB3Publisher.Enabled == nil || *config.InfluxDB3Publisher.Enabled) {
		sinks = append(sinks, data_sinks.InfluxDB3(*config.InfluxDB3Publisher))
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
		log.Fatal("No data consumers/sinks configured! Please check the config.")
	}

	log.Info("Starting processing")
	for measurement := range measurements {
		_, isOnList := filterMap[strings.ReplaceAll(measurement.Mac, ":", "")]
		if denylist && isOnList {
			log.WithFields(log.Fields{
				"mac":         measurement.Mac,
				"filter_mode": "denylist",
			}).Trace("Measurement dropped")
			continue
		}
		if allowlist && !isOnList {
			log.WithFields(log.Fields{
				"mac":         measurement.Mac,
				"filter_mode": "allowlist",
			}).Trace("Measurement dropped")
			continue
		}

		name := config.TagNames[strings.ReplaceAll(measurement.Mac, ":", "")]
		if name != "" {
			measurement.Name = &name
		} else if namedOnly {
			log.WithFields(log.Fields{
				"mac":         measurement.Mac,
				"filter_mode": "named",
			}).Trace("Measurement dropped")
			continue
		}

		if extendedValues {
			value_calculator.CalcExtendedValues(&measurement)
		}

		for _, sink := range sinks {
			sink <- measurement
		}
		log.WithFields(log.Fields{"mac": measurement.Mac}).Trace("Measurement processed")
	}
}
