package parser

import (
	"encoding/hex"
	"math"
	"testing"
)

func buildFullAdvertisementFormat5(payload []byte) []byte {
	header := []byte{0x02, 0x01, 0x04, 0x1B, 0xFF, 0x99, 0x04}
	adv := make([]byte, 0, len(header)+len(payload))
	adv = append(adv, header...)
	adv = append(adv, payload...)
	return adv
}

func TestParseFormat5_OK(t *testing.T) {
	payload := []byte{
		0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	adv := buildFullAdvertisementFormat5(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat5(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat5 returned error: %v", err)
	}

	if m.DataFormat != 0x05 {
		t.Errorf("DataFormat: got %d want %d", m.DataFormat, 0x05)
	}

	// Expected values
	expectHum := 53.49
	expectPress := 100044.0
	expectTemp := 24.3
	expectAccX := 0.004
	expectAccY := -0.004
	expectAccZ := 1.036
	expectBatt := 2.977
	expectMove := int64(66)
	expectTx := int64(4)
	expectMeas := int64(205)

	if m.Humidity == nil || int(math.Round(*m.Humidity*10000)) != int(math.Round(expectHum*10000)) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Temperature == nil || int(math.Round(*m.Temperature*1000)) != int(math.Round(expectTemp*1000)) {
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
	if m.TxPower == nil || *m.TxPower != expectTx {
		t.Errorf("TxPower: got %v want %v", m.TxPower, expectTx)
	}
	if m.MovementCounter == nil || *m.MovementCounter != expectMove {
		t.Errorf("MovementCounter: got %v want %v", m.MovementCounter, expectMove)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectMeas {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectMeas)
	}
}

func TestParseFormat5_Max(t *testing.T) {
	payload := []byte{
		0x05, 0x7F, 0xFF, 0xFF, 0xFE, 0xFF, 0xFE, 0x7F,
		0xFF, 0x7F, 0xFF, 0x7F, 0xFF, 0xFF, 0xDE, 0xFE,
		0xFF, 0xFE, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	adv := buildFullAdvertisementFormat5(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat5(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat5 returned error: %v", err)
	}

	expectHum := 163.8350
	expectPress := 115534.0
	expectTemp := 163.8350
	expectAcc := 32.767
	expectBatt := 3.646
	expectMove := int64(254)
	expectTx := int64(20)
	expectMeas := int64(65534)

	if m.Humidity == nil || int(math.Round(*m.Humidity*10000)) != int(math.Round(expectHum*10000)) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Temperature == nil || int(math.Round(*m.Temperature*1000)) != int(math.Round(expectTemp*1000)) {
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
	if m.TxPower == nil || *m.TxPower != expectTx {
		t.Errorf("TxPower: got %v want %v", m.TxPower, expectTx)
	}
	if m.MovementCounter == nil || *m.MovementCounter != expectMove {
		t.Errorf("MovementCounter: got %v want %v", m.MovementCounter, expectMove)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectMeas {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectMeas)
	}
}

func TestParseFormat5_Min(t *testing.T) {
	payload := []byte{
		0x05, 0x80, 0x01, 0x00, 0x00, 0x00, 0x00, 0x80,
		0x01, 0x80, 0x01, 0x80, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	adv := buildFullAdvertisementFormat5(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat5(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat5 returned error: %v", err)
	}

	expectHum := 0.0
	expectPress := 50000.0
	expectTemp := -163.835
	expectAcc := -32.767
	expectBatt := 1.600
	expectMove := int64(0)
	expectTx := int64(-40)
	expectMeas := int64(0)

	if m.Humidity == nil || int(math.Round(*m.Humidity*10000)) != int(math.Round(expectHum*10000)) {
		t.Errorf("Humidity: got %v want %v", m.Humidity, expectHum)
	}
	if m.Temperature == nil || int(math.Round(*m.Temperature*1000)) != int(math.Round(expectTemp*1000)) {
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
	if m.TxPower == nil || *m.TxPower != expectTx {
		t.Errorf("TxPower: got %v want %v", m.TxPower, expectTx)
	}
	if m.MovementCounter == nil || *m.MovementCounter != expectMove {
		t.Errorf("MovementCounter: got %v want %v", m.MovementCounter, expectMove)
	}
	if m.MeasurementSequenceNumber == nil || *m.MeasurementSequenceNumber != expectMeas {
		t.Errorf("MeasurementSequenceNumber: got %v want %v", m.MeasurementSequenceNumber, expectMeas)
	}
}

func TestParseFormat5_Invalid(t *testing.T) {
	payload := []byte{
		0x05, 0x80, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x80,
		0x00, 0x80, 0x00, 0x80, 0x00, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	adv := buildFullAdvertisementFormat5(payload)
	hexStr := hex.EncodeToString(adv)

	m, err := ParseFormat5(hexStr)
	if err != nil {
		t.Fatalf("ParseFormat5 returned error: %v", err)
	}

	// All fields should be nil (invalid) except DataFormat
	if m.Temperature != nil {
		t.Errorf("Temperature: expected nil, got %v", m.Temperature)
	}
	if m.Humidity != nil {
		t.Errorf("Humidity: expected nil, got %v", m.Humidity)
	}
	if m.Pressure != nil {
		t.Errorf("Pressure: expected nil, got %v", m.Pressure)
	}
	if m.AccelerationX != nil || m.AccelerationY != nil || m.AccelerationZ != nil {
		t.Errorf("Acceleration: expected nils, got X=%v Y=%v Z=%v", m.AccelerationX, m.AccelerationY, m.AccelerationZ)
	}
	if m.BatteryVoltage != nil || m.TxPower != nil {
		t.Errorf("Power fields: expected nils, got Battery=%v TxPower=%v", m.BatteryVoltage, m.TxPower)
	}
	if m.MovementCounter != nil {
		t.Errorf("MovementCounter: expected nil, got %v", m.MovementCounter)
	}
	if m.MeasurementSequenceNumber != nil {
		t.Errorf("MeasurementSequenceNumber: expected nil, got %v", m.MeasurementSequenceNumber)
	}
}
