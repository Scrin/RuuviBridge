package data_sinks

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/Scrin/RuuviBridge/common/version"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var metrics struct {
	info         prometheus.Gauge
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

	// New E1 fields
	pm10         *prometheus.GaugeVec
	pm25         *prometheus.GaugeVec
	pm40         *prometheus.GaugeVec
	pm100        *prometheus.GaugeVec
	co2          *prometheus.GaugeVec
	voc          *prometheus.GaugeVec
	nox          *prometheus.GaugeVec
	luminosity   *prometheus.GaugeVec
	soundInstant *prometheus.GaugeVec
	soundAverage *prometheus.GaugeVec
	soundPeak    *prometheus.GaugeVec

	// Diagnostics
	calibrationInProgress *prometheus.GaugeVec
	buttonPressedOnBoot   *prometheus.GaugeVec
	rtcOnBoot             *prometheus.GaugeVec
}

func initMetrics() {
	bridgeMetricPrefix := "ruuvibridge_"
	measurementMetricPrefix := "ruuvi_"
	tagLabels := []string{"name", "mac", "data_format"}

	metrics.info = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: bridgeMetricPrefix + "info",
		Help: "RuuviBridge info",
		ConstLabels: prometheus.Labels{
			"version": version.Version,
			"os":      runtime.GOOS,
			"arch":    runtime.GOARCH,
		},
	})

	metrics.measurements = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: measurementMetricPrefix + "measurements",
		Help: "Number of received measurements",
	}, tagLabels)

	metrics.temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "temperature",
		Help: "Temperature in ºC",
	}, tagLabels)
	metrics.humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "humidity",
		Help: "Relative humidity in %",
	}, tagLabels)
	metrics.pressure = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "pressure",
		Help: "Pressure in Pa",
	}, tagLabels)
	metrics.accelerationX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_x",
		Help: "X acceleration in g",
	}, tagLabels)
	metrics.accelerationY = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_y",
		Help: "Y acceleration in g",
	}, tagLabels)
	metrics.accelerationZ = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_z",
		Help: "Z acceleration in g",
	}, tagLabels)
	metrics.batteryVoltage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "battery_voltage",
		Help: "Battery voltage in V",
	}, tagLabels)
	metrics.txPower = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "tx_power",
		Help: "Transmission power in dBm",
	}, tagLabels)
	metrics.rssi = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "rssi",
		Help: "RSSI in dBm",
	}, tagLabels)
	metrics.movementCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "movement_counter",
		Help: "Number of detected movements",
	}, tagLabels)
	metrics.measurementSequenceNumber = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "measurement_sequence_number",
		Help: "Measurement sequence number",
	}, tagLabels)

	metrics.accelerationTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_total",
		Help: "Total acceleration in g",
	}, tagLabels)
	metrics.absoluteHumidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "absolute_humidity",
		Help: "Absolute humidity in g/m3",
	}, tagLabels)
	metrics.dewPoint = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "dew_point",
		Help: "Dew point in ºC",
	}, tagLabels)
	metrics.equilibriumVaporPressure = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "equilibrium_vapor_pressure",
		Help: "Equilibrium vapor pressure in Pa",
	}, tagLabels)
	metrics.airDensity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "air_density",
		Help: "Air density in kg/m3",
	}, tagLabels)
	metrics.accelerationAngleFromX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_angle_from_x",
		Help: "Acceleration angle from X in degrees",
	}, tagLabels)
	metrics.accelerationAngleFromY = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_angle_from_y",
		Help: "Acceleration angle from Y in degrees",
	}, tagLabels)
	metrics.accelerationAngleFromZ = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "acceleration_angle_from_z",
		Help: "Acceleration angle from Z in degrees",
	}, tagLabels)

	// New E1 metrics
	metrics.pm10 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "pm10",
		Help: "PM1.0 mass concentration (µg/m³)",
	}, tagLabels)
	metrics.pm25 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "pm25",
		Help: "PM2.5 mass concentration (µg/m³)",
	}, tagLabels)
	metrics.pm40 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "pm40",
		Help: "PM4.0 mass concentration (µg/m³)",
	}, tagLabels)
	metrics.pm100 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "pm100",
		Help: "PM10.0 mass concentration (µg/m³)",
	}, tagLabels)
	metrics.co2 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "co2",
		Help: "CO2 concentration (ppm)",
	}, tagLabels)
	metrics.voc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "voc",
		Help: "VOC index",
	}, tagLabels)
	metrics.nox = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "nox",
		Help: "NOx index",
	}, tagLabels)
	metrics.luminosity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "luminosity",
		Help: "Luminosity (lx)",
	}, tagLabels)
	metrics.soundInstant = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "sound_instant",
		Help: "Instant sound level (dBA)",
	}, tagLabels)
	metrics.soundAverage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "sound_average",
		Help: "Average sound level (dBA)",
	}, tagLabels)
	metrics.soundPeak = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "sound_peak",
		Help: "Peak sound level (dBA)",
	}, tagLabels)

	// Diagnostic metrics
	metrics.calibrationInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "calibration_in_progress",
		Help: "Calibration in progress (1/0)",
	}, tagLabels)
	metrics.buttonPressedOnBoot = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "button_pressed_on_boot",
		Help: "Button pressed on boot (1/0)",
	}, tagLabels)
	metrics.rtcOnBoot = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: measurementMetricPrefix + "rtc_on_boot",
		Help: "RTC was running at boot (1/0)",
	}, tagLabels)

	prometheus.MustRegister(metrics.info)
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

	// Register new E1 metrics
	prometheus.MustRegister(metrics.pm10)
	prometheus.MustRegister(metrics.pm25)
	prometheus.MustRegister(metrics.pm40)
	prometheus.MustRegister(metrics.pm100)
	prometheus.MustRegister(metrics.co2)
	prometheus.MustRegister(metrics.voc)
	prometheus.MustRegister(metrics.nox)
	prometheus.MustRegister(metrics.luminosity)
	prometheus.MustRegister(metrics.soundInstant)
	prometheus.MustRegister(metrics.soundAverage)
	prometheus.MustRegister(metrics.soundPeak)

	// Register diagnostics
	prometheus.MustRegister(metrics.calibrationInProgress)
	prometheus.MustRegister(metrics.buttonPressedOnBoot)
	prometheus.MustRegister(metrics.rtcOnBoot)

	metrics.info.Set(1)
}

func recordMetrics(m parser.Measurement) {
	name := ""
	if m.Name != nil {
		name = *m.Name
	}
	labels := prometheus.Labels{"name": name, "mac": m.Mac, "data_format": fmt.Sprintf("%X", m.DataFormat)}
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
	safeSetB := func(gauge *prometheus.GaugeVec, v *bool) {
		if v != nil {
			if *v {
				gauge.With(labels).Set(1)
			} else {
				gauge.With(labels).Set(0)
			}
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

	// New E1 fields
	safeSetF(metrics.pm10, m.Pm10)
	safeSetF(metrics.pm25, m.Pm25)
	safeSetF(metrics.pm40, m.Pm40)
	safeSetF(metrics.pm100, m.Pm100)
	safeSetF(metrics.co2, m.CO2)
	safeSetF(metrics.voc, m.VOC)
	safeSetF(metrics.nox, m.NOX)
	safeSetF(metrics.luminosity, m.Illuminance)
	safeSetF(metrics.soundInstant, m.SoundInstant)
	safeSetF(metrics.soundAverage, m.SoundAverage)
	safeSetF(metrics.soundPeak, m.SoundPeak)

	// Diagnostics
	safeSetB(metrics.calibrationInProgress, m.CalibrationInProgress)
	safeSetB(metrics.buttonPressedOnBoot, m.ButtonPressedOnBoot)
	safeSetB(metrics.rtcOnBoot, m.RtcOnBoot)
}

func Prometheus(conf config.Prometheus) chan<- parser.Measurement {
	port := conf.Port
	if port == 0 {
		port = 8081
	}
	log.WithFields(log.Fields{"port": port}).Info("Starting prometheus sink")
	measurements := make(chan parser.Measurement, 1024)
	initMetrics()
	go func() {
		for measurement := range measurements {
			recordMetrics(measurement)
		}
	}()

	go http.ListenAndServe(fmt.Sprintf(":%d", port), promhttp.Handler())

	return measurements
}
