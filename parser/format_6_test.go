package parser

import (
	"encoding/hex"
	"math"
	"testing"
)

func buildFullAdvertisementFormat6(payload []byte) []byte {
	// Build proper BLE advertisement structure based on parser expectations
	// Bytes 0-7: BLE header, Byte 8: 0xFF, Bytes 9-10: Company ID, Bytes 11+: Payload
	adv := make([]byte, 11+len(payload))
	// BLE advertisement header (8 bytes)
	adv[0] = 0x2B // Length field or similar
	adv[1] = 0x01 // AD Type
	adv[2] = 0x06 // Flags
	adv[3] = 0x00 // Padding
	// Manufacturer specific data
	adv[4] = 0xFF // Manufacturer specific data flag
	adv[5] = 0x99 // Company identifier (Ruuvi)
	adv[6] = 0x04 // Company identifier (Ruuvi)
	// Copy payload
	copy(adv[7:], payload)
	return adv
}

func floatEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestParseFormat6_ValidData(t *testing.T) {
	// Case: valid data
	// Raw binary data: 0x06170C5668C79E007000C90501D9XXCD004C884F (XX = reserved)
	payload := []byte{
		0x06,       // Data format
		0x17, 0x0C, // Temperature (29.500 C)
		0x56, 0x68, // Humidity (55.300%)
		0xC7, 0x9E, // Pressure (101102 Pa)
		0x00, 0x70, // PM2.5 (11.2 ug/m^3)
		0x00, 0xC9, // CO2 (201 ppm)
		0x05,             // VOC (10)
		0x01,             // NOX (2)
		0xD9,             // Luminosity (13026.67 Lux)
		0xFF,             // Reserved
		0xCD,             // Measurement Sequence (205)
		0x00,             // Flags
		0x4C, 0x88, 0x4F, // MAC (4C 88 4F)
	}

	adv := buildFullAdvertisementFormat6(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat6(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat6 returned error: %v", err)
	}

	// Verify data format
	if m.DataFormat != 0x06 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0x06)
	}

	// Verify temperature: 29.500 C
	expectedTemp := 29.500
	if m.Temperature == nil || *m.Temperature != expectedTemp {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectedTemp)
	}

	// Verify humidity: 55.300%
	expectedHumidity := 55.300
	if m.Humidity == nil || !floatEqual(*m.Humidity, expectedHumidity, 0.001) {
		t.Errorf("Humidity: got %v want %v", *m.Humidity, expectedHumidity)
	}

	// Verify pressure: 101102 Pa
	expectedPressure := 101102.0
	if m.Pressure == nil || *m.Pressure != expectedPressure {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectedPressure)
	}

	// Verify PM2.5: 11.2 ug/m^3
	expectedPM25 := 11.2
	if m.Pm2p5 == nil || *m.Pm2p5 != expectedPM25 {
		t.Errorf("PM2.5: got %v want %v", m.Pm2p5, expectedPM25)
	}

	// Verify CO2: 201 ppm
	expectedCO2 := 201.0
	if m.CO2 == nil || *m.CO2 != expectedCO2 {
		t.Errorf("CO2: got %v want %v", m.CO2, expectedCO2)
	}

	// Verify VOC: 10
	expectedVOC := 10.0
	if m.VOC == nil || *m.VOC != expectedVOC {
		t.Errorf("VOC: got %v want %v", m.VOC, expectedVOC)
	}

	// Verify NOX: 2
	expectedNOX := 2.0
	if m.NOX == nil || *m.NOX != expectedNOX {
		t.Errorf("NOX: got %v want %v", m.NOX, expectedNOX)
	}

	// Verify luminosity: 13026.67 Lux (this will help us verify the formula)
	expectedLuminosity := 13026.67
	if m.Illuminance == nil {
		t.Errorf("Illuminance: got nil want %v", expectedLuminosity)
	} else if !floatEqual(*m.Illuminance, expectedLuminosity, 0.01) {
		t.Errorf("Illuminance: got %v want %v (diff: %v)", *m.Illuminance, expectedLuminosity, *m.Illuminance-expectedLuminosity)
	}

	// Verify measurement sequence: 205
	expectedSeq := int64(205)
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectedSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectedSeq)
	}
}

func TestParseFormat6_MaximumValues(t *testing.T) {
	// Case: maximum values
	// Raw binary data: 0x067FFF9C40FFFE27109C40FAFAFEXXFF074C8F4F (XX = reserved)
	payload := []byte{
		0x06,       // Data format
		0x7F, 0xFF, // Temperature (163.835 C)
		0x9C, 0x40, // Humidity (100.000%)
		0xFF, 0xFE, // Pressure (115534 Pa)
		0x27, 0x10, // PM2.5 (1000.0 ug/m^3)
		0x9C, 0x40, // CO2 (40000 ppm)
		0xFA,             // VOC (500)
		0xFA,             // NOX (500)
		0xFE,             // Luminosity (65355.00 Lux)
		0xFF,             // Reserved
		0xFF,             // Measurement Sequence (255)
		0x07,             // Flags (calibration in progress)
		0x4C, 0x8F, 0x4F, // MAC
	}

	adv := buildFullAdvertisementFormat6(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat6(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat6 returned error: %v", err)
	}

	// Verify temperature: 163.835 C
	expectedTemp := 163.835
	if m.Temperature == nil || *m.Temperature != expectedTemp {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectedTemp)
	}

	// Verify humidity: 100.000%
	expectedHumidity := 100.000
	if m.Humidity == nil || *m.Humidity != expectedHumidity {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectedHumidity)
	}

	// Verify pressure: 115534 Pa
	expectedPressure := 115534.0
	if m.Pressure == nil || *m.Pressure != expectedPressure {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectedPressure)
	}

	// Verify PM2.5: 1000.0 ug/m^3
	expectedPM25 := 1000.0
	if m.Pm2p5 == nil || *m.Pm2p5 != expectedPM25 {
		t.Errorf("PM2.5: got %v want %v", m.Pm2p5, expectedPM25)
	}

	// Verify CO2: 40000 ppm
	expectedCO2 := 40000.0
	if m.CO2 == nil || *m.CO2 != expectedCO2 {
		t.Errorf("CO2: got %v want %v", m.CO2, expectedCO2)
	}

	// Verify VOC: 500
	expectedVOC := 500.0
	if m.VOC == nil || *m.VOC != expectedVOC {
		t.Errorf("VOC: got %v want %v", m.VOC, expectedVOC)
	}

	// Verify NOX: 500
	expectedNOX := 500.0
	if m.NOX == nil || *m.NOX != expectedNOX {
		t.Errorf("NOX: got %v want %v", m.NOX, expectedNOX)
	}

	// Verify luminosity: 65535.00 Lux (0xFE according to documentation)
	expectedLuminosity := 65535.00
	if m.Illuminance == nil {
		t.Errorf("Illuminance: got nil want %v", expectedLuminosity)
	} else if !floatEqual(*m.Illuminance, expectedLuminosity, 1.0) {
		t.Errorf("Illuminance: got %v want %v (diff: %v)", *m.Illuminance, expectedLuminosity, *m.Illuminance-expectedLuminosity)
	}

	// Verify measurement sequence: 255
	expectedSeq := int64(255)
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectedSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectedSeq)
	}
}

func TestParseFormat6_MinimumValues(t *testing.T) {
	// Case: minimum values
	// Raw binary data: 0x0680010000000000000000000000XX00004C884F (XX = reserved)
	payload := []byte{
		0x06,       // Data format
		0x80, 0x01, // Temperature (-163.835 C)
		0x00, 0x00, // Humidity (0.000%)
		0x00, 0x00, // Pressure (50000 Pa)
		0x00, 0x00, // PM2.5 (0.0 ug/m^3)
		0x00, 0x00, // CO2 (0 ppm)
		0x00,             // VOC (0)
		0x00,             // NOX (0)
		0x00,             // Luminosity (0.00 Lux)
		0xFF,             // Reserved
		0x00,             // Measurement Sequence (0)
		0x00,             // Flags
		0x4C, 0x88, 0x4F, // MAC
	}

	adv := buildFullAdvertisementFormat6(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat6(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat6 returned error: %v", err)
	}

	// Verify temperature: -163.835 C
	expectedTemp := -163.835
	if m.Temperature == nil || *m.Temperature != expectedTemp {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectedTemp)
	}

	// Verify humidity: 0.000%
	expectedHumidity := 0.000
	if m.Humidity == nil || *m.Humidity != expectedHumidity {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectedHumidity)
	}

	// Verify pressure: 50000 Pa
	expectedPressure := 50000.0
	if m.Pressure == nil || *m.Pressure != expectedPressure {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectedPressure)
	}

	// Verify PM2.5: 0.0 ug/m^3
	expectedPM25 := 0.0
	if m.Pm2p5 == nil || *m.Pm2p5 != expectedPM25 {
		t.Errorf("PM2.5: got %v want %v", m.Pm2p5, expectedPM25)
	}

	// Verify CO2: 0 ppm
	expectedCO2 := 0.0
	if m.CO2 == nil || *m.CO2 != expectedCO2 {
		t.Errorf("CO2: got %v want %v", m.CO2, expectedCO2)
	}

	// Verify VOC: 0
	expectedVOC := 0.0
	if m.VOC == nil || *m.VOC != expectedVOC {
		t.Errorf("VOC: got %v want %v", m.VOC, expectedVOC)
	}

	// Verify NOX: 0
	expectedNOX := 0.0
	if m.NOX == nil || *m.NOX != expectedNOX {
		t.Errorf("NOX: got %v want %v", m.NOX, expectedNOX)
	}

	// Verify luminosity: 0.00 Lux
	expectedLuminosity := 0.0
	if m.Illuminance == nil || *m.Illuminance != expectedLuminosity {
		t.Errorf("Illuminance: got %v want %v", m.Illuminance, expectedLuminosity)
	}

	// Verify measurement sequence: 0
	expectedSeq := int64(0)
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectedSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectedSeq)
	}
}

func TestParseFormat6_InvalidValues(t *testing.T) {
	// Case: Invalid values (all 0xFF means no data available)
	// Raw binary data: 0x068000FFFFFFFFFFFFFFFFFFFFFFFFXXFFFFFFFFFF (XX = reserved)
	payload := []byte{
		0x06,       // Data format
		0x80, 0x00, // Temperature (invalid - should be ignored)
		0xFF, 0xFF, // Humidity (invalid)
		0xFF, 0xFF, // Pressure (invalid)
		0xFF, 0xFF, // PM2.5 (invalid)
		0xFF, 0xFF, // CO2 (invalid)
		0xFF,             // VOC (invalid)
		0xFF,             // NOX (invalid)
		0xFF,             // Luminosity (invalid)
		0xFF,             // Reserved
		0xFF,             // Measurement Sequence (255)
		0xFF,             // Flags (all set)
		0xFF, 0xFF, 0xFF, // MAC (invalid)
	}

	adv := buildFullAdvertisementFormat6(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat6(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat6 returned error: %v", err)
	}

	// Verify data format
	if m.DataFormat != 0x06 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0x06)
	}

	// Temperature should be nil (invalid value 0x8000)
	if m.Temperature != nil {
		t.Errorf("Temperature: got %v want nil (invalid value)", *m.Temperature)
	}

	// All other fields should be nil due to invalid values (0xFFFF)
	if m.Humidity != nil {
		t.Errorf("Humidity: got %v want nil (invalid value)", *m.Humidity)
	}
	if m.Pressure != nil {
		t.Errorf("Pressure: got %v want nil (invalid value)", *m.Pressure)
	}
	if m.Pm2p5 != nil {
		t.Errorf("PM2.5: got %v want nil (invalid value)", *m.Pm2p5)
	}
	if m.CO2 != nil {
		t.Errorf("CO2: got %v want nil (invalid value)", *m.CO2)
	}
	if m.VOC != nil {
		t.Errorf("VOC: got %v want nil (invalid value)", *m.VOC)
	}
	if m.NOX != nil {
		t.Errorf("NOX: got %v want nil (invalid value)", *m.NOX)
	}
	if m.Illuminance != nil {
		t.Errorf("Illuminance: got %v want nil (invalid value)", *m.Illuminance)
	}

	// Measurement sequence should be 255 (valid value)
	expectedSeq := int64(255)
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectedSeq {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectedSeq)
	}
}
