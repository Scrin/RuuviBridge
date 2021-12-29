package data_sinks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
)

func InfluxDB(conf config.InfluxDBPublisher) chan<- parser.Measurement {
	url := conf.Url
	if url == "" {
		url = "https://localhost:8086"
	}
	bucket := conf.Bucket
	if bucket == "" {
		bucket = "ruuvi"
	}
	measurementName := conf.Measurement
	if measurementName == "" {
		measurementName = "ruuvi_measurements"
	}
	log.WithFields(log.Fields{"target": url}).Info("Starting InfluxDB sink")

	client := influxdb.NewClient(url, conf.AuthToken)
	writeAPI := client.WriteAPIBlocking(conf.Org, bucket)

	measurements := make(chan parser.Measurement)
	go func() {
		for measurement := range measurements {
			p := influxdb.NewPointWithMeasurement(measurementName).
				AddTag("dataFormat", fmt.Sprintf("%d", measurement.DataFormat)).
				AddTag("mac", strings.ReplaceAll(measurement.Mac, ":", ""))
			if measurement.Name != nil {
				p.AddTag("name", *measurement.Name)
			}
			addFloat(p, "temperature", measurement.Temperature)
			addFloat(p, "humidity", measurement.Humidity)
			addFloat(p, "pressure", measurement.Pressure)
			addFloat(p, "accelerationX", measurement.AccelerationX)
			addFloat(p, "accelerationY", measurement.AccelerationY)
			addFloat(p, "accelerationZ", measurement.AccelerationZ)
			addFloat(p, "batteryVoltage", measurement.BatteryVoltage)
			addInt(p, "txPower", measurement.TxPower)
			addInt(p, "rssi", measurement.Rssi)
			addInt(p, "movementCounter", measurement.MovementCounter)
			addInt(p, "measurementSequenceNumber", measurement.MeasurementSequenceNumber)
			addFloat(p, "accelerationTotal", measurement.AccelerationTotal)
			addFloat(p, "absoluteHumidity", measurement.AbsoluteHumidity)
			addFloat(p, "dewPoint", measurement.DewPoint)
			addFloat(p, "equilibriumVaporPressure", measurement.EquilibriumVaporPressure)
			addFloat(p, "airDensity", measurement.AirDensity)
			addFloat(p, "accelerationAngleFromX", measurement.AccelerationAngleFromX)
			addFloat(p, "accelerationAngleFromY", measurement.AccelerationAngleFromY)
			addFloat(p, "accelerationAngleFromZ", measurement.AccelerationAngleFromZ)
			p.SetTime(time.Now())
			err := writeAPI.WritePoint(context.Background(), p)
			if err != nil {
				log.Error("Failed to send data to InfluxDB: ", err)
			}
		}
		client.Close()
	}()
	return measurements
}

func addFloat(p *write.Point, name string, value *float64) {
	if value != nil {
		p.AddField(name, *value)
	}
}

func addInt(p *write.Point, name string, value *int64) {
	if value != nil {
		p.AddField(name, *value)
	}
}
