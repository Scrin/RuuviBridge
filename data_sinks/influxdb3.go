package data_sinks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/Scrin/RuuviBridge/common/limiter"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	log "github.com/sirupsen/logrus"
)

func InfluxDB3(conf config.InfluxDB3Publisher) chan<- parser.Measurement {
	url := conf.Url
	if url == "" {
		url = "https://localhost:8086"
	}
	measurementName := conf.Measurement
	if measurementName == "" {
		measurementName = "ruuvi_measurements"
	}
	log.WithFields(log.Fields{
		"target":           url,
		"measurement_name": measurementName,
		"minimum_interval": conf.MinimumInterval,
	}).Info("Starting InfluxDB3 sink")

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     url,
		Token:    conf.AuthToken,
		Database: conf.Database,
	})

	if err != nil {
		log.WithError(err).Error("Failed to create InfluxDB3 client")
	}

	limiter := limiter.New(conf.MinimumInterval)
	measurements := make(chan parser.Measurement, 1024)
	go func() {
		for measurement := range measurements {
			if !limiter.Check(measurement) {
				log.WithFields(log.Fields{"mac": measurement.Mac}).Trace("Skipping InfluxDB3 publish due to interval limit")
				continue
			}
			go func(measurement parser.Measurement) {
				p := influxdb3.NewPointWithMeasurement(measurementName).
					AddTag("dataFormat", fmt.Sprintf("%d", measurement.DataFormat)).
					AddTag("mac", strings.ReplaceAll(measurement.Mac, ":", ""))
				if measurement.Name != nil {
					p.AddTag("name", *measurement.Name)
				}
				for tag, value := range conf.AdditionalTags {
					p.AddTag(tag, value)
				}
				influx3AddFloat(p, "temperature", measurement.Temperature)
				influx3AddFloat(p, "humidity", measurement.Humidity)
				influx3AddFloat(p, "pressure", measurement.Pressure)
				influx3AddFloat(p, "accelerationX", measurement.AccelerationX)
				influx3AddFloat(p, "accelerationY", measurement.AccelerationY)
				influx3AddFloat(p, "accelerationZ", measurement.AccelerationZ)
				influx3AddFloat(p, "batteryVoltage", measurement.BatteryVoltage)
				influx3AddInt(p, "txPower", measurement.TxPower)
				influx3AddInt(p, "rssi", measurement.Rssi)
				influx3AddInt(p, "movementCounter", measurement.MovementCounter)
				influx3AddInt(p, "measurementSequenceNumber", measurement.MeasurementSequenceNumber)
				influx3AddFloat(p, "accelerationTotal", measurement.AccelerationTotal)
				influx3AddFloat(p, "absoluteHumidity", measurement.AbsoluteHumidity)
				influx3AddFloat(p, "dewPoint", measurement.DewPoint)
				influx3AddFloat(p, "equilibriumVaporPressure", measurement.EquilibriumVaporPressure)
				influx3AddFloat(p, "airDensity", measurement.AirDensity)
				influx3AddFloat(p, "accelerationAngleFromX", measurement.AccelerationAngleFromX)
				influx3AddFloat(p, "accelerationAngleFromY", measurement.AccelerationAngleFromY)
				influx3AddFloat(p, "accelerationAngleFromZ", measurement.AccelerationAngleFromZ)
				p.SetTimestamp(time.Now())
				err := client.WritePoints(context.Background(), p)
				if err != nil {
					log.WithError(err).Error("Failed to send data to InfluxDB3")
				}
			}(measurement)
		}
		client.Close()
	}()
	return measurements
}

func influx3AddFloat(p *influxdb3.Point, name string, value *float64) {
	if value != nil {
		p.AddField(name, *value)
	}
}

func influx3AddInt(p *influxdb3.Point, name string, value *int64) {
	if value != nil {
		p.AddField(name, *value)
	}
}
