package parser

import (
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
		return m, errors.New("Data is too short")
	}

	if data[4] != 0xff { // manufacturer specific data
		return m, errors.New("Data is not manufacturer specific data")
	}

	if data[5] != ruuviCompanyIdentifier[0] || data[6] != ruuviCompanyIdentifier[1] {
		return m, errors.New("Data has wrong company identifier")
	}

	data = data[7:]

	if data[0] != 5 { // data format
		return m, errors.New("Data is not in data format 5")
	}

	m.DataFormat = int64(data[0])
	m.Temperature = f64(float64(int16(binary.BigEndian.Uint16(data[1:]))) / 200)
	m.Humidity = f64(float64(binary.BigEndian.Uint16(data[3:])) / 400)
	m.Pressure = f64(float64(binary.BigEndian.Uint16(data[5:])) + 50_000)
	m.AccelerationX = f64(float64(int16(binary.BigEndian.Uint16(data[7:]))) / 1000)
	m.AccelerationY = f64(float64(int16(binary.BigEndian.Uint16(data[9:]))) / 1000)
	m.AccelerationZ = f64(float64(int16(binary.BigEndian.Uint16(data[11:]))) / 1000)

	powerInfo := binary.BigEndian.Uint16(data[13:])

	m.BatteryVoltage = f64(float64(powerInfo>>5)/1000 + 1.6)
	m.TxPower = i64(int64(powerInfo&0b11111)*2 - 40)

	m.MovementCounter = i64(int64(data[15]))
	m.MeasurementSequenceNumber = i64(int64(binary.BigEndian.Uint16(data[16:])))

	log.WithFields(log.Fields{
		"raw_data":    input,
		"data_format": m.DataFormat,
	}).Trace("Successfully parsed data")
	return m, nil
}
