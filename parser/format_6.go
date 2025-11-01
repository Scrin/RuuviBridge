package parser

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"

	"github.com/rs/zerolog/log"
)

func ParseFormat6(input string) (Measurement, error) {
	var m Measurement
	data, err := hex.DecodeString(input)
	if err != nil {
		return m, err
	}
	if len(data) < 23 {
		return m, errors.New("data is too short")
	}

	if data[4] != 0xff { // manufacturer specific data
		return m, errors.New("data is not manufacturer specific data")
	}

	if data[5] != ruuviCompanyIdentifier[0] || data[6] != ruuviCompanyIdentifier[1] {
		return m, errors.New("data has wrong company identifier")
	}

	data = data[7:]

	if data[0] != 0x06 { // data format
		return m, errors.New("data is not in data format 6")
	}

	m.DataFormat = int64(data[0])

	// Temperature (offset 1-2): -32767 ... 32767, 0.005 degrees resolution
	if !bytes.Equal(data[1:3], []byte{0x80, 0x00}) {
		m.Temperature = f64(float64(int16(binary.BigEndian.Uint16(data[1:3]))) * 0.005)
	}

	// Humidity (offset 3-4): 0 ... 40000, 0.0025% resolution
	if !bytes.Equal(data[3:5], []byte{0xff, 0xff}) {
		m.Humidity = f64(float64(binary.BigEndian.Uint16(data[3:5])) * 0.0025)
	}

	// Pressure (offset 5-6): 0 ... 65534, 1 Pa units with -50000 Pa offset
	if !bytes.Equal(data[5:7], []byte{0xff, 0xff}) {
		m.Pressure = f64(float64(binary.BigEndian.Uint16(data[5:7])) + 50000)
	}

	// PM2.5 (offset 7-8): 0 ... 10000, 0.1 resolution
	if !bytes.Equal(data[7:9], []byte{0xff, 0xff}) {
		m.Pm2p5 = f64(float64(binary.BigEndian.Uint16(data[7:9])) / 10)
	}

	// CO2 (offset 9-10): 0 ... 40000, 1 resolution
	if !bytes.Equal(data[9:11], []byte{0xff, 0xff}) {
		m.CO2 = f64(float64(binary.BigEndian.Uint16(data[9:11])))
	}

	// Flags byte (offset 16)
	flags := data[16]
	voc9 := (flags >> 6) & 0x01
	nox9 := (flags >> 7) & 0x01

	// VOC (offset 11 + FLAGS b6): 9-bit unsigned, 1 resolution
	combinedVOC := ((uint16(data[11]) << 1) | uint16(voc9))
	if combinedVOC != 0x1FF {
		m.VOC = f64(float64(combinedVOC))
	}

	// NOX (offset 12 + FLAGS b7): 9-bit unsigned, 1 resolution
	combinedNOX := ((uint16(data[12]) << 1) | uint16(nox9))
	if combinedNOX != 0x1FF {
		m.NOX = f64(float64(combinedNOX))
	}

	// Luminosity (offset 13): Logarithmic, range 0 ... 65535
	if data[13] != 0xff {
		// Convert logarithmic value to linear lux using official formula:
		// MAX_VALUE := 65535, MAX_CODE := 254
		// DELTA := ln(MAX_VALUE + 1) / MAX_CODE
		// VALUE := exp(CODE * DELTA) - 1
		const maxValue = 65535.0
		const maxCode = 254.0
		delta := math.Log(maxValue+1) / maxCode
		lumValue := math.Exp(float64(data[13])*delta) - 1
		m.Illuminance = f64(lumValue)
	}

	// Measurement sequence (offset 15): LSB of full sequence counter
	// 255 is a valid value, so we always set it
	m.MeasurementSequenceNumber = i64(int64(data[15]))

	log.Trace().
		Str("raw_data", input).
		Int64("data_format", m.DataFormat).
		Msg("Successfully parsed data")
	return m, nil
}
