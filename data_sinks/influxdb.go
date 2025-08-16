package data_sinks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Scrin/RuuviBridge/common/limiter"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/rs/zerolog/log"
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
	log.Info().
		Str("target", url).
		Str("bucket", bucket).
		Str("measurement_name", measurementName).
		Dur("minimum_interval", conf.MinimumInterval).
		Msg("Starting InfluxDB sink")

	client := influxdb.NewClient(url, conf.AuthToken)
	writeAPI := client.WriteAPIBlocking(conf.Org, bucket)

	limiter := limiter.New(conf.MinimumInterval)
	measurements := make(chan parser.Measurement, 1024)
	go func() {
		for measurement := range measurements {
			if !limiter.Check(measurement) {
				log.Trace().Str("mac", measurement.Mac).Msg("Skipping InfluxDB publish due to interval limit")
				continue
			}
			go func(measurement parser.Measurement) {
				p := influxdb.NewPointWithMeasurement(measurementName).
					AddTag("dataFormat", fmt.Sprintf("%X", measurement.DataFormat)).
					AddTag("mac", strings.ReplaceAll(measurement.Mac, ":", ""))
				if measurement.Name != nil {
					p.AddTag("name", *measurement.Name)
				}
				for tag, value := range conf.AdditionalTags {
					p.AddTag(tag, value)
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
				// New E1 fields
				addFloat(p, "pm10", measurement.Pm10)
				addFloat(p, "pm25", measurement.Pm25)
				addFloat(p, "pm40", measurement.Pm40)
				addFloat(p, "pm100", measurement.Pm100)
				addFloat(p, "co2", measurement.CO2)
				addFloat(p, "voc", measurement.VOC)
				addFloat(p, "nox", measurement.NOX)
				addFloat(p, "illuminance", measurement.Illuminance)
				addFloat(p, "soundInstant", measurement.SoundInstant)
				addFloat(p, "soundAverage", measurement.SoundAverage)
				addFloat(p, "soundPeak", measurement.SoundPeak)
				// Diagnostics
				addBool(p, "calibrationInProgress", measurement.CalibrationInProgress)
				addBool(p, "buttonPressedOnBoot", measurement.ButtonPressedOnBoot)
				addBool(p, "rtcOnBoot", measurement.RtcOnBoot)
				p.SetTime(time.Now())
				err := writeAPI.WritePoint(context.Background(), p)
				if err != nil {
					log.Error().Err(err).Msg("Failed to send data to InfluxDB")
				}
			}(measurement)
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

func addBool(p *write.Point, name string, value *bool) {
	if value != nil {
		p.AddField(name, *value)
	}
}
