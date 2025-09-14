package parser

type Measurement struct {
	CommonData
	BasicEnvironmentalData
	AirQualityData
	DiagnosticsData
	UnofficialData
	CalculatedData
}

// Common data for all measurements
type CommonData struct {
	Name       *string `json:"name,omitempty"`
	Mac        string  `json:"mac,omitempty"`
	Timestamp  *int64  `json:"timestamp,omitempty"`
	DataFormat int64   `json:"data_format,omitempty"`
}

// Basic environmental data, typically on ruuvitags
type BasicEnvironmentalData struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	Humidity        *float64 `json:"humidity,omitempty"`
	Pressure        *float64 `json:"pressure,omitempty"`
	AccelerationX   *float64 `json:"accelerationX,omitempty"`
	AccelerationY   *float64 `json:"accelerationY,omitempty"`
	AccelerationZ   *float64 `json:"accelerationZ,omitempty"`
	BatteryVoltage  *float64 `json:"batteryVoltage,omitempty"`
	TxPower         *int64   `json:"txPower,omitempty"`
	Rssi            *int64   `json:"rssi,omitempty"`
	MovementCounter *int64   `json:"movementCounter,omitempty"`
}

// Air quality data, typically on <redacted>
type AirQualityData struct {
	Pm1p0       *float64 `json:"pm1p0,omitempty"`
	Pm2p5       *float64 `json:"pm2p5,omitempty"`
	Pm4p0       *float64 `json:"pm4p0,omitempty"`
	Pm10p0      *float64 `json:"pm10p0,omitempty"`
	CO2         *float64 `json:"co2,omitempty"`
	VOC         *float64 `json:"voc,omitempty"`
	NOX         *float64 `json:"nox,omitempty"`
	Illuminance *float64 `json:"illuminance,omitempty"`
}

// Diagnostics data
type DiagnosticsData struct {
	MeasurementSequenceNumber *int64 `json:"measurementSequenceNumber,omitempty"`
	CalibrationInProgress     *bool  `json:"calibrationInProgress,omitempty"`
}

// Data not officially documented (eg. on format E1, transmitted by certain revisions of <redacted>)
type UnofficialData struct {
	SoundInstant        *float64 `json:"soundInstant,omitempty"`
	SoundAverage        *float64 `json:"soundAverage,omitempty"`
	SoundPeak           *float64 `json:"soundPeak,omitempty"`
	ButtonPressedOnBoot *bool    `json:"buttonPressedOnBoot,omitempty"`
	RtcOnBoot           *bool    `json:"rtcOnBoot,omitempty"`
}

// Calculated data not actually present on measurements, but instead calculated
type CalculatedData struct {
	AccelerationTotal        *float64 `json:"accelerationTotal,omitempty"`
	AbsoluteHumidity         *float64 `json:"absoluteHumidity,omitempty"`
	DewPoint                 *float64 `json:"dewPoint,omitempty"`
	EquilibriumVaporPressure *float64 `json:"equilibriumVaporPressure,omitempty"`
	AirDensity               *float64 `json:"airDensity,omitempty"`
	AccelerationAngleFromX   *float64 `json:"accelerationAngleFromX,omitempty"`
	AccelerationAngleFromY   *float64 `json:"accelerationAngleFromY,omitempty"`
	AccelerationAngleFromZ   *float64 `json:"accelerationAngleFromZ,omitempty"`
	AirQualityIndex          *float64 `json:"airQualityIndex,omitempty"`
}
