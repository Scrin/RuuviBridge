package parser

import (
	"encoding/hex"
	"math"
	"testing"
)

func buildFullAdvertisement(payload []byte) []byte {
	header := []byte{0x2B, 0xFF, 0x99, 0x04}
	adv := make([]byte, 0, len(header)+len(payload))
	adv = append(adv, header...)
	adv = append(adv, payload...)
	return adv
}

func roundInt10(v float64) int {
	return int(math.Round(v * 10.0))
}

func TestParseFormatE1_OK(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x17, 0x0C, // Temperature (29.5 C -> 5900/200)
		0x56, 0x68, // Humidity (55.3 % -> 22120/400)
		0xC7, 0x9E, // Pressure (101102 Pa -> 51102 + 50000)
		0x00, 0x65, // PM1.0 (10.1)
		0x00, 0x70, // PM2.5 (11.2)
		0x04, 0xBD, // PM4.0 (121.3)
		0x11, 0xCA, // PM10.0 (455.4)
		0x00, 0xC9, // CO2 (201)
		0x05,             // VOC LSB (combined -> 10)
		0x01,             // NOX LSB (combined -> 2)
		0x13, 0xE0, 0xAC, // Luminosity (13027.00 -> uint24 / 100)
		0x3D,             // Sound inst LSB (-> 42.4)
		0x4A,             // Sound avg LSB (-> 47.6)
		0x9C,             // Sound peak LSB (-> 80.4)
		0xDE, 0xCD, 0xEE, // Seq cnt (0xDECDEE)
		0x00,                         // Flags (all false, 9th bits 0)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	// Expected values from the firmware test
	expectTemp := 29.5
	expectHum := 55.3
	expectPress := 101102.0
	expectPM10 := 10.1
	expectPM25 := 11.2
	expectPM40 := 121.3
	expectPM100 := 455.4
	expectCO2 := 201.0
	expectVOC := 10.0
	expectNOX := 2.0
	expectLux := 13027.0
	expectSoundInst := 42.4
	expectSoundAvg := 47.6
	expectSoundPeak := 80.4
	expectSeq := int64(0x00DECDEE)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_Zeroes(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_Temperature(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x13, 0x88, // Temperature (25.0 C -> 5000/200)
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 25.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_Humidity(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x6D, 0x60, // Humidity (70.0 % -> 28000/400)
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 70.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_Pressure(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0xC3, 0x50, // Pressure (100000 Pa -> 50000 + 50000)
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 100000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_PM1p0(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x1B, 0x58, // PM1.0 (700.0)
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 700.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_PM2p5(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x1B, 0x58, // PM2.5 (700.0)
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 700.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_PM4p0(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x1B, 0x58, // PM4.0 (700.0)
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 700.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_PM10p0(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x1B, 0x58, // PM10.0 (700.0)
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 700.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_CO2(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x88, 0xB8, // CO2 (35000)
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 35000.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_VOC(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0xF9,             // VOC (499, 9th bit set)
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x40,                         // Flags (VOC 9th bit)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 499.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_NOX(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0xF8,             // NOX (497, 9th bit set)
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x80,                         // Flags (NOX 9th bit)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 497.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_Luminosity(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0xDB, 0xBA, 0x02, // Luminosity (144000.02)
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 144000.02
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_SoundInstant(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0xCF,             // Sound inst (101.0, 9th bit set)
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x08,                         // Flags (Sound inst 9th bit)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 101.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_SoundAverage(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0xCF,             // Sound avg (101.0, 9th bit set)
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x10,                         // Flags (Sound avg 9th bit)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 101.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_SoundPeak(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0xCF,             // Sound peak (101.0, 9th bit set)
		0x00, 0x00, 0x00, // Seq cnt
		0x20,                         // Flags (Sound peak 9th bit)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 101.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_SeqCnt(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0xAB, 0xCD, 0xEF, // Seq cnt (0xABCDEF)
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 0.0
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0x00ABCDEF)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}

func TestParseFormatE1_FlagCalibrationInProgress(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x01,                         // Flags (calibration)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectCal := true
	expectBtn := false
	expectRtc := false
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != expectCal {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, expectCal)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != expectBtn {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, expectBtn)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != expectRtc {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, expectRtc)
	}
}

func TestParseFormatE1_FlagButtonPressed(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x02,                         // Flags (button)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectCal := false
	expectBtn := true
	expectRtc := false
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != expectCal {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, expectCal)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != expectBtn {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, expectBtn)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != expectRtc {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, expectRtc)
	}
}

func TestParseFormatE1_FlagRtcRunningOnBoot(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x00, 0x00, // Temperature
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x04,                         // Flags (RTC)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectCal := false
	expectBtn := false
	expectRtc := true
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != expectCal {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, expectCal)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != expectBtn {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, expectBtn)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != expectRtc {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, expectRtc)
	}
}

func TestParseFormatE1_Max(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x7F, 0xFF, // Temperature (163.835)
		0x9C, 0x40, // Humidity (100.0)
		0xFF, 0xFE, // Pressure (115534)
		0x27, 0x10, // PM1.0 (1000.0)
		0x27, 0x10, // PM2.5 (1000.0)
		0x27, 0x10, // PM4.0 (1000.0)
		0x27, 0x10, // PM10.0 (1000.0)
		0x9C, 0x40, // CO2 (40000)
		0xFA,             // VOC (500)
		0xFA,             // NOX (500)
		0xDC, 0x28, 0xF0, // Luminosity (144284.00)
		0xFF,             // Sound inst (120)
		0xFF,             // Sound avg (120)
		0xFF,             // Sound peak (120)
		0xFF, 0xFF, 0xFE, // Seq cnt (0xFFFFFE)
		0x07,                         // Flags (all true)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := 163.835
	expectHum := 100.0
	expectPress := 115534.0
	expectPM10 := 1000.0
	expectPM25 := 1000.0
	expectPM40 := 1000.0
	expectPM100 := 1000.0
	expectCO2 := 40000.0
	expectVOC := 500.0
	expectNOX := 500.0
	expectLux := 144284.00
	expectSoundInst := 120.0
	expectSoundAvg := 120.0
	expectSoundPeak := 120.0
	expectSeq := int64(0x00FFFFFE)
	expectCal := true
	expectBtn := true
	expectRtc := true

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != expectCal {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, expectCal)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != expectBtn {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, expectBtn)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != expectRtc {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, expectRtc)
	}
}

func TestParseFormatE1_Min(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x80, 0x01, // Temperature (-163.835)
		0x00, 0x00, // Humidity
		0x00, 0x00, // Pressure
		0x00, 0x00, // PM1.0
		0x00, 0x00, // PM2.5
		0x00, 0x00, // PM4.0
		0x00, 0x00, // PM10.0
		0x00, 0x00, // CO2
		0x00,             // VOC
		0x00,             // NOX
		0x00, 0x00, 0x00, // Luminosity
		0x00,             // Sound inst
		0x00,             // Sound avg
		0x00,             // Sound peak
		0x00, 0x00, 0x00, // Seq cnt
		0x00,                         // Flags
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	expectTemp := -163.835
	expectHum := 0.0
	expectPress := 50000.0
	expectPM10 := 0.0
	expectPM25 := 0.0
	expectPM40 := 0.0
	expectPM100 := 0.0
	expectCO2 := 0.0
	expectVOC := 0.0
	expectNOX := 0.0
	expectLux := 0.0
	expectSoundInst := 18.0
	expectSoundAvg := 18.0
	expectSoundPeak := 18.0
	expectSeq := int64(0)

	if m.DataFormat != 0xE1 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0xE1)
	}
	if m.Temperature == nil || roundInt10(*m.Temperature) != roundInt10(expectTemp) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Humidity == nil || roundInt10(*m.Humidity) != roundInt10(expectHum) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.Pm10 == nil || roundInt10(*m.Pm10) != roundInt10(expectPM10) {
		t.Errorf("Pm10: got %v want %v", m.Pm10, expectPM10)
	}
	if m.Pm25 == nil || roundInt10(*m.Pm25) != roundInt10(expectPM25) {
		t.Errorf("Pm25: got %v want %v", m.Pm25, expectPM25)
	}
	if m.Pm40 == nil || roundInt10(*m.Pm40) != roundInt10(expectPM40) {
		t.Errorf("Pm40: got %v want %v", m.Pm40, expectPM40)
	}
	if m.Pm100 == nil || roundInt10(*m.Pm100) != roundInt10(expectPM100) {
		t.Errorf("Pm100: got %v want %v", m.Pm100, expectPM100)
	}
	if m.CO2 == nil || int(math.Round(*m.CO2)) != int(math.Round(expectCO2)) {
		t.Errorf("CO2: got %v want %v", m.CO2, expectCO2)
	}
	if m.VOC == nil || int(math.Round(*m.VOC)) != int(math.Round(expectVOC)) {
		t.Errorf("VOC: got %v want %v", m.VOC, expectVOC)
	}
	if m.NOX == nil || int(math.Round(*m.NOX)) != int(math.Round(expectNOX)) {
		t.Errorf("NOX: got %v want %v", m.NOX, expectNOX)
	}
	if m.Illuminance == nil || int(math.Round(*m.Illuminance)) != int(math.Round(expectLux)) {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectLux)
	}
	if m.SoundInstant == nil || roundInt10(*m.SoundInstant) != roundInt10(expectSoundInst) {
		t.Errorf("SoundInstant: got %v want %v", m.SoundInstant, expectSoundInst)
	}
	if m.SoundAverage == nil || roundInt10(*m.SoundAverage) != roundInt10(expectSoundAvg) {
		t.Errorf("SoundAverage: got %v want %v", m.SoundAverage, expectSoundAvg)
	}
	if m.SoundPeak == nil || roundInt10(*m.SoundPeak) != roundInt10(expectSoundPeak) {
		t.Errorf("SoundPeak: got %v want %v", m.SoundPeak, expectSoundPeak)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectSeq)
	}
}

func TestParseFormatE1_InvalidData(t *testing.T) {
	payload := []byte{
		0xE1,       // Data type
		0x80, 0x00, // Temperature (invalid)
		0xFF, 0xFF, // Humidity (invalid)
		0xFF, 0xFF, // Pressure (invalid)
		0xFF, 0xFF, // PM1.0 (invalid)
		0xFF, 0xFF, // PM2.5 (invalid)
		0xFF, 0xFF, // PM4.0 (invalid)
		0xFF, 0xFF, // PM10.0 (invalid)
		0xFF, 0xFF, // CO2 (invalid)
		0xFF,             // VOC LSB (invalid)
		0xFF,             // NOX LSB (invalid)
		0xFF, 0xFF, 0xFF, // Luminosity (invalid)
		0xFF,             // Sound inst LSB
		0xFF,             // Sound avg LSB
		0xFF,             // Sound peak LSB
		0xFF, 0xFF, 0xFF, // Seq cnt (invalid)
		0xF8,                         // Flags (diag all true; VOC/NOX 9th=1; sound 9th=0)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Reserved
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // MAC address
	}
	adv := buildFullAdvertisement(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormatE1(hexStr)
	if err != nil {
		t.Fatalf("ParseFormatE1 returned error: %v", err)
	}

	// Invalid fields should be nil
	if m.Temperature != nil {
		t.Errorf("Temperature: expected nil, got %v", m.Temperature)
	}
	if m.Humidity != nil {
		t.Errorf("Humidity: expected nil, got %v", m.Humidity)
	}
	if m.Pressure != nil {
		t.Errorf("Pressure: expected nil, got %v", m.Pressure)
	}
	if m.Pm10 != nil {
		t.Errorf("Pm10: expected nil, got %v", m.Pm10)
	}
	if m.Pm25 != nil {
		t.Errorf("Pm25: expected nil, got %v", m.Pm25)
	}
	if m.Pm40 != nil {
		t.Errorf("Pm40: expected nil, got %v", m.Pm40)
	}
	if m.Pm100 != nil {
		t.Errorf("Pm100: expected nil, got %v", m.Pm100)
	}
	if m.CO2 != nil {
		t.Errorf("CO2: expected nil, got %v", m.CO2)
	}
	if m.VOC != nil {
		t.Errorf("VOC: expected nil, got %v", m.VOC)
	}
	if m.NOX != nil {
		t.Errorf("NOX: expected nil, got %v", m.NOX)
	}
	if m.Illuminance != nil {
		t.Errorf("Illuminance: expected nil, got %v", m.Illuminance)
	}
	if m.MeasurementSequenceNumber != nil {
		t.Errorf("MeasurementSequenceNumber: expected nil, got %v", m.MeasurementSequenceNumber)
	}

	// Sound values should be considered invalid and thus nil with these bytes
	if m.SoundInstant != nil {
		t.Errorf("SoundInstant: expected nil, got %v", m.SoundInstant)
	}
	if m.SoundAverage != nil {
		t.Errorf("SoundAverage: expected nil, got %v", m.SoundAverage)
	}
	if m.SoundPeak != nil {
		t.Errorf("SoundPeak: expected nil, got %v", m.SoundPeak)
	}

	// Flags: all diagnostic flags false per 0xF8 (bits 0..2 are 0)
	if m.CalibrationInProgress == nil || *m.CalibrationInProgress != false {
		t.Errorf("CalibrationInProgress: got %v want %v", m.CalibrationInProgress, false)
	}
	if m.ButtonPressedOnBoot == nil || *m.ButtonPressedOnBoot != false {
		t.Errorf("ButtonPressedOnBoot: got %v want %v", m.ButtonPressedOnBoot, false)
	}
	if m.RtcOnBoot == nil || *m.RtcOnBoot != false {
		t.Errorf("RtcOnBoot: got %v want %v", m.RtcOnBoot, false)
	}
}
