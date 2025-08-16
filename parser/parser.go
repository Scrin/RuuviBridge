package parser

import "github.com/rs/zerolog/log"

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

	Pm10         *float64 `json:"pm10,omitempty"`
	Pm25         *float64 `json:"pm25,omitempty"`
	Pm40         *float64 `json:"pm40,omitempty"`
	Pm100        *float64 `json:"pm100,omitempty"`
	CO2          *float64 `json:"co2,omitempty"`
	VOC          *float64 `json:"voc,omitempty"`
	NOX          *float64 `json:"nox,omitempty"`
	Illuminance  *float64 `json:"illuminance,omitempty"`
	SoundInstant *float64 `json:"soundInstant,omitempty"`
	SoundAverage *float64 `json:"soundAverage,omitempty"`
	SoundPeak    *float64 `json:"soundPeak,omitempty"`

	CalibrationInProgress *bool `json:"calibrationInProgress,omitempty"`
	ButtonPressedOnBoot   *bool `json:"buttonPressedOnBoot,omitempty"`
	RtcOnBoot             *bool `json:"rtcOnBoot,omitempty"`

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
	var err_formate1, err_format5, err_format3 error
	if measurement, err_formate1 = ParseFormatE1(input); err_formate1 == nil {
		return measurement, true
	}
	if measurement, err_format5 = ParseFormat5(input); err_format5 == nil {
		return measurement, true
	}
	if measurement, err_format3 = ParseFormat3(input); err_format3 == nil {
		return measurement, true
	}
	log.Trace().
		Str("raw_data", input).
		Str("format_e1_error", err_formate1.Error()).
		Str("format_5_error", err_format5.Error()).
		Str("format_3_error", err_format3.Error()).
		Msg("Failed to parse data")
	return Measurement{}, false
}
