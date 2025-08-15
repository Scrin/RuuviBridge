package parser

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"

	log "github.com/sirupsen/logrus"
)

func ParseFormat5(input string) (Measurement, error) {
	var m Measurement
	data, err := hex.DecodeString(input)
	if err != nil {
		return m, err
	}
	if len(data) < 31 {
		return m, errors.New("data is too short")
	}

	if data[4] != 0xff { // manufacturer specific data
		return m, errors.New("data is not manufacturer specific data")
	}

	if data[5] != ruuviCompanyIdentifier[0] || data[6] != ruuviCompanyIdentifier[1] {
		return m, errors.New("data has wrong company identifier")
	}

	data = data[7:]

	if data[0] != 0x05 { // data format
		return m, errors.New("data is not in data format 5")
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
	if !bytes.Equal(data[7:9], []byte{0x80, 0x00}) {
		m.AccelerationX = f64(float64(int16(binary.BigEndian.Uint16(data[7:9]))) / 1000)
	}
	if !bytes.Equal(data[9:11], []byte{0x80, 0x00}) {
		m.AccelerationY = f64(float64(int16(binary.BigEndian.Uint16(data[9:11]))) / 1000)
	}
	if !bytes.Equal(data[11:13], []byte{0x80, 0x00}) {
		m.AccelerationZ = f64(float64(int16(binary.BigEndian.Uint16(data[11:13]))) / 1000)
	}
	if !bytes.Equal(data[13:15], []byte{0xff, 0xff}) {
		powerInfo := binary.BigEndian.Uint16(data[13:15])
		m.BatteryVoltage = f64(float64(powerInfo>>5)/1000 + 1.6)
		m.TxPower = i64(int64(powerInfo&0b11111)*2 - 40)
	}
	if data[15] != 0xff {
		m.MovementCounter = i64(int64(data[15]))
	}
	if !bytes.Equal(data[16:18], []byte{0xff, 0xff}) {
		m.MeasurementSequenceNumber = i64(int64(binary.BigEndian.Uint16(data[16:18])))
	}

	log.WithFields(log.Fields{
		"raw_data":    input,
		"data_format": m.DataFormat,
	}).Trace("Successfully parsed data")
	return m, nil
}
