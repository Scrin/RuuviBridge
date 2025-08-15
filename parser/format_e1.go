package parser

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"

	log "github.com/sirupsen/logrus"
)

func ParseFormatE1(input string) (Measurement, error) {
	var m Measurement
	data, err := hex.DecodeString(input)
	if err != nil {
		return m, err
	}
	if len(data) < 32 {
		return m, errors.New("data is too short")
	}

	if data[1] != 0xff { // manufacturer specific data
		return m, errors.New("data is not manufacturer specific data")
	}

	if data[2] != ruuviCompanyIdentifier[0] || data[3] != ruuviCompanyIdentifier[1] {
		return m, errors.New("data has wrong company identifier")
	}

	data = data[4:]

	if data[0] != 0xe1 { // data format
		return m, errors.New("data is not in data format E1")
	}

	m.DataFormat = int64(data[0])
	if !bytes.Equal(data[1:3], []byte{0x80, 0x00}) {
		m.Temperature = f64(float64(int16(binary.BigEndian.Uint16(data[1:3]))) / 200)
	}
	if !bytes.Equal(data[3:5], []byte{0xff, 0xff}) {
		m.Humidity = f64(float64(binary.BigEndian.Uint16(data[3:5])) / 400)
	}
	if !bytes.Equal(data[5:7], []byte{0xff, 0xff}) {
		m.Pressure = f64(float64(binary.BigEndian.Uint16(data[5:7])) + 50_000)
	}

	// Particulate matter (0.1 resolution)
	if !bytes.Equal(data[7:9], []byte{0xff, 0xff}) {
		m.Pm10 = f64(float64(binary.BigEndian.Uint16(data[7:9])) / 10)
	}
	if !bytes.Equal(data[9:11], []byte{0xff, 0xff}) {
		m.Pm25 = f64(float64(binary.BigEndian.Uint16(data[9:11])) / 10)
	}
	if !bytes.Equal(data[11:13], []byte{0xff, 0xff}) {
		m.Pm40 = f64(float64(binary.BigEndian.Uint16(data[11:13])) / 10)
	}
	if !bytes.Equal(data[13:15], []byte{0xff, 0xff}) {
		m.Pm100 = f64(float64(binary.BigEndian.Uint16(data[13:15])) / 10)
	}

	// CO2 (1 resolution)
	if !bytes.Equal(data[15:17], []byte{0xff, 0xff}) {
		m.CO2 = f64(float64(binary.BigEndian.Uint16(data[15:17])))
	}

	// VOC and NOX (now 9-bit unsigned using special flags for MSB)
	// LSB 8 bits in data[17] and data[18], 9th bit from flags at data[28]

	// Luminosity (uint24, 0.01 resolution)
	if !bytes.Equal(data[19:22], []byte{0xff, 0xff, 0xff}) {
		lux := (uint32(data[19]) << 16) | (uint32(data[20]) << 8) | uint32(data[21])
		m.Illuminance = f64(float64(lux) / 100)
	}

	// Sound levels (now 9-bit unsigned using special flags for MSB), 0.2 resolution with +18 offset
	// LSB 8 bits in data[22..24], 9th bit from flags at data[28]

	// Measurement sequence number (uint24)
	if !bytes.Equal(data[25:28], []byte{0xff, 0xff, 0xff}) {
		seq := (uint32(data[25]) << 16) | (uint32(data[26]) << 8) | uint32(data[27])
		m.MeasurementSequenceNumber = i64(int64(seq))
	}

	// Special flags bitfield at data[28]
	flags := data[28]
	// Bits 0..2: diagnostic flags (calibration, button, rtc)
	calInProgress := (flags & 0x01) != 0
	buttonPressed := (flags & 0x02) != 0
	rtcOnBoot := (flags & 0x04) != 0
	m.CalibrationInProgress = &calInProgress
	m.ButtonPressedOnBoot = &buttonPressed
	m.RtcOnBoot = &rtcOnBoot
	// Bits 3..7: dbaInst9, dbaAvg9, soundPeak9, voc9, nox9 respectively
	dbaInst9 := (flags >> 3) & 0x01
	dbaAvg9 := (flags >> 4) & 0x01
	soundPeak9 := (flags >> 5) & 0x01
	voc9 := (flags >> 6) & 0x01
	nox9 := (flags >> 7) & 0x01

	// Upgrade sound values to 9-bit
	combinedSoundInstant := ((uint16(data[22]) << 1) | uint16(dbaInst9))
	combinedSoundAverage := ((uint16(data[23]) << 1) | uint16(dbaAvg9))
	combinedSoundPeak := ((uint16(data[24]) << 1) | uint16(soundPeak9))
	if combinedSoundInstant != 0x1FF {
		m.SoundInstant = f64(float64(combinedSoundInstant)*0.2 + 18)
	}
	if combinedSoundAverage != 0x1FF {
		m.SoundAverage = f64(float64(combinedSoundAverage)*0.2 + 18)
	}
	if combinedSoundPeak != 0x1FF {
		m.SoundPeak = f64(float64(combinedSoundPeak)*0.2 + 18)
	}

	// Upgrade VOC/NOX to 9-bit unsigned
	combinedVOC := ((uint16(data[17]) << 1) | uint16(voc9))
	combinedNOX := ((uint16(data[18]) << 1) | uint16(nox9))
	if combinedVOC != 0x1FF {
		m.VOC = f64(float64(combinedVOC))
	}
	if combinedNOX != 0x1FF {
		m.NOX = f64(float64(combinedNOX))
	}

	log.WithFields(log.Fields{
		"raw_data":    input,
		"data_format": m.DataFormat,
	}).Trace("Successfully parsed data")
	return m, nil
}
