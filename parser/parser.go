package parser

import log "github.com/sirupsen/logrus"

type Measurement struct {
	Name       *string `json:"name,omitempty"`
	Mac        string  `json:"mac,omitempty"`
	Timestamp  *int64  `json:"timestamp,omitempty"`
	DataFormat int64   `json:"data_format,omitempty"`

	Temperature               *float64 `json:"temperature,omitempty"`
	Humidity                  *float64 `json:"humidity,omitempty"`
	Pressure                  *float64 `json:"pressure,omitempty"`
	AccelerationX             *float64 `json:"accelerationX,omitempty"`
	AccelerationY             *float64 `json:"accelerationY,omitempty"`
	AccelerationZ             *float64 `json:"accelerationZ,omitempty"`
	BatteryVoltage            *float64 `json:"batteryVoltage,omitempty"`
	TxPower                   *int64   `json:"txPower,omitempty"`
	Rssi                      *int64   `json:"rssi,omitempty"`
	MovementCounter           *int64   `json:"movementCounter,omitempty"`
	MeasurementSequenceNumber *int64   `json:"measurementSequenceNumber,omitempty"`

	AccelerationTotal        *float64 `json:"accelerationTotal,omitempty"`
	AbsoluteHumidity         *float64 `json:"absoluteHumidity,omitempty"`
	DewPoint                 *float64 `json:"dewPoint,omitempty"`
	EquilibriumVaporPressure *float64 `json:"equilibriumVaporPressure,omitempty"`
	AirDensity               *float64 `json:"airDensity,omitempty"`
	AccelerationAngleFromX   *float64 `json:"accelerationAngleFromX,omitempty"`
	AccelerationAngleFromY   *float64 `json:"accelerationAngleFromY,omitempty"`
	AccelerationAngleFromZ   *float64 `json:"accelerationAngleFromZ,omitempty"`
}

var ruuviCompanyIdentifier = []byte{0x99, 0x04} // 0x0499

func f64(value float64) *float64 {
	return &value
}
func i64(value int64) *int64 {
	return &value
}

func Parse(input string) (Measurement, bool) {
	var measurement Measurement
	var err_format5, err_format3 error
	if measurement, err_format5 = ParseFormat5(input); err_format5 == nil {
		return measurement, true
	}
	if measurement, err_format3 = ParseFormat3(input); err_format3 == nil {
		return measurement, true
	}
	log.WithFields(log.Fields{
		"raw_data":      input,
		"format5_error": err_format5,
		"format3_error": err_format3,
	}).Trace("Failed to parse data")
	return Measurement{}, false
}
