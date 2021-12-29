package data_sinks

import (
	"fmt"
	"net/http"

	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var metrics struct {
	measurements *prometheus.CounterVec

	temperature               *prometheus.GaugeVec
	humidity                  *prometheus.GaugeVec
	pressure                  *prometheus.GaugeVec
	accelerationX             *prometheus.GaugeVec
	accelerationY             *prometheus.GaugeVec
	accelerationZ             *prometheus.GaugeVec
	batteryVoltage            *prometheus.GaugeVec
	txPower                   *prometheus.GaugeVec
	rssi                      *prometheus.GaugeVec
	movementCounter           *prometheus.GaugeVec
	measurementSequenceNumber *prometheus.GaugeVec

	accelerationTotal        *prometheus.GaugeVec
	absoluteHumidity         *prometheus.GaugeVec
	dewPoint                 *prometheus.GaugeVec
	equilibriumVaporPressure *prometheus.GaugeVec
	airDensity               *prometheus.GaugeVec
	accelerationAngleFromX   *prometheus.GaugeVec
	accelerationAngleFromY   *prometheus.GaugeVec
	accelerationAngleFromZ   *prometheus.GaugeVec
}

func initMetrics() {
	metricPrefix := "ruuvitag_"
	tagLabels := []string{"name", "mac", "data_format"}

	metrics.measurements = prometheus.NewCounterVec(prometheus.CounterOpts{Name: metricPrefix + "measurements", Help: "Number of received measurements"}, tagLabels)

	metrics.temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "temperature", Help: "Temperature in ºC"}, tagLabels)
	metrics.humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "humidity", Help: "Relative humidity in %"}, tagLabels)
	metrics.pressure = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "pressure", Help: "Pressure in Pa"}, tagLabels)
	metrics.accelerationX = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_x", Help: "X acceleration in g"}, tagLabels)
	metrics.accelerationY = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_y", Help: "Y acceleration in g"}, tagLabels)
	metrics.accelerationZ = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_z", Help: "Z acceleration in g"}, tagLabels)
	metrics.batteryVoltage = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "battery_voltage", Help: "Battery voltage in V"}, tagLabels)
	metrics.txPower = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "tx_power", Help: "Transmission power in dBm"}, tagLabels)
	metrics.rssi = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "rssi", Help: "RSSI in dBm"}, tagLabels)
	metrics.movementCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "movement_counter", Help: "Number of detected movements"}, tagLabels)
	metrics.measurementSequenceNumber = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "measurement_sequence_number", Help: "Measurement sequence number"}, tagLabels)

	metrics.accelerationTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_total", Help: "Total acceleration in g"}, tagLabels)
	metrics.absoluteHumidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "absolute_humidity", Help: "Absolute humidity in g/m3"}, tagLabels)
	metrics.dewPoint = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "dew_point", Help: "Dew point in ºC"}, tagLabels)
	metrics.equilibriumVaporPressure = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "equilibrium_vapor_pressure", Help: "Equilibrium vapor pressure in Pa"}, tagLabels)
	metrics.airDensity = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "air_density", Help: "Air density in kg/m3"}, tagLabels)
	metrics.accelerationAngleFromX = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_angle_from_x", Help: "Acceleration angle from X in degrees"}, tagLabels)
	metrics.accelerationAngleFromY = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_angle_from_y", Help: "Acceleration angle from Y in degrees"}, tagLabels)
	metrics.accelerationAngleFromZ = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "acceleration_angle_from_z", Help: "Acceleration angle from Z in degrees"}, tagLabels)

	prometheus.MustRegister(metrics.measurements)

	prometheus.MustRegister(metrics.temperature)
	prometheus.MustRegister(metrics.humidity)
	prometheus.MustRegister(metrics.pressure)
	prometheus.MustRegister(metrics.accelerationX)
	prometheus.MustRegister(metrics.accelerationY)
	prometheus.MustRegister(metrics.accelerationZ)
	prometheus.MustRegister(metrics.batteryVoltage)
	prometheus.MustRegister(metrics.txPower)
	prometheus.MustRegister(metrics.rssi)
	prometheus.MustRegister(metrics.movementCounter)
	prometheus.MustRegister(metrics.measurementSequenceNumber)

	prometheus.MustRegister(metrics.accelerationTotal)
	prometheus.MustRegister(metrics.absoluteHumidity)
	prometheus.MustRegister(metrics.dewPoint)
	prometheus.MustRegister(metrics.equilibriumVaporPressure)
	prometheus.MustRegister(metrics.airDensity)
	prometheus.MustRegister(metrics.accelerationAngleFromX)
	prometheus.MustRegister(metrics.accelerationAngleFromY)
	prometheus.MustRegister(metrics.accelerationAngleFromZ)
}

func recordMetrics(m parser.Measurement) {
	name := ""
	if m.Name != nil {
		name = *m.Name
	}
	labels := prometheus.Labels{"name": name, "mac": m.Mac, "data_format": fmt.Sprint(m.DataFormat)}
	safeSetF := func(gauge *prometheus.GaugeVec, v *float64) {
		if v != nil {
			gauge.With(labels).Set(*v)
		}
	}
	safeSetI := func(gauge *prometheus.GaugeVec, v *int64) {
		if v != nil {
			gauge.With(labels).Set(float64(*v))
		}
	}

	metrics.measurements.With(labels).Inc()

	safeSetF(metrics.temperature, m.Temperature)
	safeSetF(metrics.humidity, m.Humidity)
	safeSetF(metrics.pressure, m.Pressure)
	safeSetF(metrics.accelerationX, m.AccelerationX)
	safeSetF(metrics.accelerationY, m.AccelerationY)
	safeSetF(metrics.accelerationZ, m.AccelerationZ)
	safeSetF(metrics.batteryVoltage, m.BatteryVoltage)
	safeSetI(metrics.txPower, m.TxPower)
	safeSetI(metrics.rssi, m.Rssi)
	safeSetI(metrics.movementCounter, m.MovementCounter)
	safeSetI(metrics.measurementSequenceNumber, m.MeasurementSequenceNumber)

	safeSetF(metrics.accelerationTotal, m.AccelerationTotal)
	safeSetF(metrics.absoluteHumidity, m.AbsoluteHumidity)
	safeSetF(metrics.dewPoint, m.DewPoint)
	safeSetF(metrics.equilibriumVaporPressure, m.EquilibriumVaporPressure)
	safeSetF(metrics.airDensity, m.AirDensity)
	safeSetF(metrics.accelerationAngleFromX, m.AccelerationAngleFromX)
	safeSetF(metrics.accelerationAngleFromY, m.AccelerationAngleFromY)
	safeSetF(metrics.accelerationAngleFromZ, m.AccelerationAngleFromZ)
}

func Prometheus(conf config.Prometheus) chan<- parser.Measurement {
	port := conf.Port
	if port == 0 {
		port = 8080
	}
	log.Info("Starting prometheus sink")
	measurements := make(chan parser.Measurement)
	initMetrics()
	go func() {
		for measurement := range measurements {
			recordMetrics(measurement)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return measurements
}
