package parser

import (
	"encoding/hex"
	"math"
	"testing"
)

func buildFullAdvertisementFormat3(payload []byte) []byte {
	header := []byte{0x02, 0x01, 0x04, 0x1B, 0xFF, 0x99, 0x04}
	adv := make([]byte, 0, len(header)+len(payload))
	adv = append(adv, header...)
	adv = append(adv, payload...)
	return adv
}

func TestParseFormat3_OK(t *testing.T) {
	payload := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC, 0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	adv := buildFullAdvertisementFormat3(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat3(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat3 returned error: %v", err)
	}

	if m.DataFormat != 0x03 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0x03)
	}

	expectHum := 20.5
	expectPress := 102766.0
	expectTemp := 26.3
	expectAccX := -1.000
	expectAccY := -1.726
	expectAccZ := 0.714
	expectBatt := 2.899

	if m.Humidity == nil || int(math.Round(*m.Humidity*10)) != int(math.Round(expectHum*10)) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Temperature == nil || int(math.Round(*m.Temperature*100)) != int(math.Round(expectTemp*100)) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.AccelerationX == nil || int(math.Round(*m.AccelerationX*1000)) != int(math.Round(expectAccX*1000)) {
		t.Errorf("AccelerationX: got %v want %v", m.AccelerationX, expectAccX)
	}
	if m.AccelerationY == nil || int(math.Round(*m.AccelerationY*1000)) != int(math.Round(expectAccY*1000)) {
		t.Errorf("AccelerationY: got %v want %v", m.AccelerationY, expectAccY)
	}
	if m.AccelerationZ == nil || int(math.Round(*m.AccelerationZ*1000)) != int(math.Round(expectAccZ*1000)) {
		t.Errorf("AccelerationZ: got %v want %v", m.AccelerationZ, expectAccZ)
	}
	if m.BatteryVoltage == nil || int(math.Round(*m.BatteryVoltage*1000)) != int(math.Round(expectBatt*1000)) {
		t.Errorf("BatteryVoltage: got %v want %v", m.BatteryVoltage, expectBatt)
	}
}

func TestParseFormat3_Max(t *testing.T) {
	payload := []byte{
		0x03, 0xFF, 0x7F, 0x63, 0xFF, 0xFF, 0x7F, 0xFF, 0x7F, 0xFF, 0x7F, 0xFF, 0xFF, 0xFF,
	}
	adv := buildFullAdvertisementFormat3(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat3(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat3 returned error: %v", err)
	}

	expectHum := 127.5
	expectPress := 115535.0
	expectTemp := 127.99
	expectAcc := 32.767
	expectBatt := 65.535

	if m.Humidity == nil || int(math.Round(*m.Humidity*10)) != int(math.Round(expectHum*10)) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Temperature == nil || int(math.Round(*m.Temperature*100)) != int(math.Round(expectTemp*100)) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.AccelerationX == nil || int(math.Round(*m.AccelerationX*1000)) != int(math.Round(expectAcc*1000)) {
		t.Errorf("AccelerationX: got %v want %v", m.AccelerationX, expectAcc)
	}
	if m.AccelerationY == nil || int(math.Round(*m.AccelerationY*1000)) != int(math.Round(expectAcc*1000)) {
		t.Errorf("AccelerationY: got %v want %v", m.AccelerationY, expectAcc)
	}
	if m.AccelerationZ == nil || int(math.Round(*m.AccelerationZ*1000)) != int(math.Round(expectAcc*1000)) {
		t.Errorf("AccelerationZ: got %v want %v", m.AccelerationZ, expectAcc)
	}
	if m.BatteryVoltage == nil || int(math.Round(*m.BatteryVoltage*1000)) != int(math.Round(expectBatt*1000)) {
		t.Errorf("BatteryVoltage: got %v want %v", m.BatteryVoltage, expectBatt)
	}
}

func TestParseFormat3_Min(t *testing.T) {
	payload := []byte{
		0x03, 0x00, 0xFF, 0x63, 0x00, 0x00, 0x80, 0x01, 0x80, 0x01, 0x80, 0x01, 0x00, 0x00,
	}
	adv := buildFullAdvertisementFormat3(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat3(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat3 returned error: %v", err)
	}

	expectHum := 0.0
	expectPress := 50000.0
	expectTemp := -127.99
	expectAcc := -32.767
	expectBatt := 0.000

	if m.Humidity == nil || int(math.Round(*m.Humidity*10)) != int(math.Round(expectHum*10)) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Temperature == nil || int(math.Round(*m.Temperature*100)) != int(math.Round(expectTemp*100)) {
		t.Errorf("Temperature: got %v want %v", m.Temperature, expectTemp)
	}
	if m.Pressure == nil || int(math.Round(*m.Pressure)) != int(math.Round(expectPress)) {
		t.Errorf("Pressure: got %v want %v", m.Pressure, expectPress)
	}
	if m.AccelerationX == nil || int(math.Round(*m.AccelerationX*1000)) != int(math.Round(expectAcc*1000)) {
		t.Errorf("AccelerationX: got %v want %v", m.AccelerationX, expectAcc)
	}
	if m.AccelerationY == nil || int(math.Round(*m.AccelerationY*1000)) != int(math.Round(expectAcc*1000)) {
		t.Errorf("AccelerationY: got %v want %v", m.AccelerationY, expectAcc)
	}
	if m.AccelerationZ == nil || int(math.Round(*m.AccelerationZ*1000)) != int(math.Round(expectAcc*1000)) {
		t.Errorf("AccelerationZ: got %v want %v", m.AccelerationZ, expectAcc)
	}
	if m.BatteryVoltage == nil || int(math.Round(*m.BatteryVoltage*1000)) != int(math.Round(expectBatt*1000)) {
		t.Errorf("BatteryVoltage: got %v want %v", m.BatteryVoltage, expectBatt)
	}
}
