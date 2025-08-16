package parser

import (
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/rs/zerolog/log"
)

func ParseFormat3(input string) (Measurement, error) {
	var m Measurement
	data, err := hex.DecodeString(input)
	if err != nil {
		return m, err
	}
	if len(data) < 21 {
		return m, errors.New("data is too short")
	}

	if data[4] != 0xff { // manufacturer specific data
		return m, errors.New("data is not manufacturer specific data")
	}

	if data[5] != ruuviCompanyIdentifier[0] || data[6] != ruuviCompanyIdentifier[1] {
		return m, errors.New("data has wrong company identifier")
	}

	data = data[7:]

	if data[0] != 0x03 { // data format
		return m, errors.New("data is not in data format 3")
	}

	m.DataFormat = int64(data[0])
	m.Humidity = f64(float64(data[1]) / 2)
	temperatureSign := (data[2] >> 7) & 1
	temperatureBase := data[2] & 0x7F
	temperatureFraction := float64(data[3]) / 100
	temperature := float64(temperatureBase) + temperatureFraction
	if temperatureSign == 1 {
		temperature *= -1
	}
	m.Temperature = f64(temperature)
	m.Pressure = f64(float64(binary.BigEndian.Uint16(data[4:])) + 50_000)
	m.AccelerationX = f64(float64(int16(binary.BigEndian.Uint16(data[6:]))) / 1000)
	m.AccelerationY = f64(float64(int16(binary.BigEndian.Uint16(data[8:]))) / 1000)
	m.AccelerationZ = f64(float64(int16(binary.BigEndian.Uint16(data[10:]))) / 1000)
	m.BatteryVoltage = f64(float64(binary.BigEndian.Uint16(data[12:])) / 1000)

	log.Trace().
		Str("raw_data", input).
		Int64("data_format", m.DataFormat).
		Msg("Successfully parsed data")
	return m, nil
}
