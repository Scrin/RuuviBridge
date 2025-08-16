package value_calculator

import (
	"math"

	"github.com/Scrin/RuuviBridge/parser"
)

func CalcExtendedValues(m *parser.Measurement) {
	// from https://github.com/Scrin/RuuviCollector/blob/master/src/main/java/fi/tkgwf/ruuvi/utils/MeasurementValueCalculator.java
	f64 := func(value float64) *float64 { return &value }
	if m.AccelerationX != nil && m.AccelerationY != nil && m.AccelerationZ != nil {
		m.AccelerationTotal = f64(math.Sqrt((*m.AccelerationX)*(*m.AccelerationX) + (*m.AccelerationY)*(*m.AccelerationY) + (*m.AccelerationZ)*(*m.AccelerationZ)))
	}
	if m.AccelerationX != nil && m.AccelerationTotal != nil && *m.AccelerationTotal != 0 {
		m.AccelerationAngleFromX = f64(math.Acos((*m.AccelerationX)/(*m.AccelerationTotal)) * (180 / math.Pi))
	}
	if m.AccelerationY != nil && m.AccelerationTotal != nil && *m.AccelerationTotal != 0 {
		m.AccelerationAngleFromY = f64(math.Acos((*m.AccelerationY)/(*m.AccelerationTotal)) * (180 / math.Pi))
	}
	if m.AccelerationZ != nil && m.AccelerationTotal != nil && *m.AccelerationTotal != 0 {
		m.AccelerationAngleFromZ = f64(math.Acos((*m.AccelerationZ)/(*m.AccelerationTotal)) * (180 / math.Pi))
	}
	if m.Temperature != nil {
		m.EquilibriumVaporPressure = f64(611.2 * math.Exp(17.67*(*m.Temperature)/(243.5+(*m.Temperature))))
	}
	if m.Temperature != nil && m.Humidity != nil {
		m.AbsoluteHumidity = f64((*m.EquilibriumVaporPressure) * (*m.Humidity) * 0.021674 / (273.15 + (*m.Temperature)))
	}
	if m.EquilibriumVaporPressure != nil && m.Humidity != nil && *m.Humidity != 0 {
		v := math.Log((*m.Humidity) / 100 * (*m.EquilibriumVaporPressure) / 611.2)
		m.DewPoint = f64(-243.5 * v / (v - 17.67))
	}
	if m.Temperature != nil && m.Humidity != nil && m.Pressure != nil && m.EquilibriumVaporPressure != nil {
		m.AirDensity = f64(1.2929 * 273.15 / (*m.Temperature + 273.15) * (float64(*m.Pressure) - 0.3783*(*m.Humidity)/100*(*m.EquilibriumVaporPressure)) / 101300)
	}
	if m.Pm2p5 != nil && m.CO2 != nil {
		const aqiMax = 100.0
		const pm25Min = 0.0
		const pm25Max = 60.0
		const co2Min = 420.0
		const co2Max = 2300.0
		const pm25Scale = aqiMax / (pm25Max - pm25Min)
		const co2Scale = aqiMax / (co2Max - co2Min)

		pm25 := math.Max(pm25Min, math.Min(pm25Max, *m.Pm2p5))
		co2 := math.Max(co2Min, math.Min(co2Max, *m.CO2))

		dx := (pm25 - pm25Min) * pm25Scale
		dy := (co2 - co2Min) * co2Scale

		r := math.Hypot(dx, dy)
		aqi := math.Max(0, math.Min(aqiMax, aqiMax-r))
		m.AirQualityIndex = f64(aqi)
	}
}
