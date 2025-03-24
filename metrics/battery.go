package metrics

import (
	"machine"
	"time"
)

var batteryPin machine.ADC

const (
	nSamples      = 128
	adcResolution = 12
	minBatteryV   = 3.0
	maxBatteryV   = 4.2
)

func InitBattery() {
	machine.InitADC()
	batteryPin = machine.ADC{Pin: machine.P0_04}
	batteryPin.Configure(machine.ADCConfig{
		Reference:  minBatteryV,
		Resolution: adcResolution,
	})
	time.Sleep(1 * time.Second)
}

func ReadBatteryLevel() float32 {
	return readBatteryVoltage()
}

func readBatteryVoltage() float32 {
	time.Sleep(1 * time.Second)
	raw := readAverageADC(batteryPin, nSamples)
	batteryV := (float32(raw) / float32(4095) * minBatteryV
	if batteryV < minBatteryV {
		batteryV = minBatteryV
	} else if batteryV > maxBatteryV {
		batteryV = maxBatteryV
	}
	return batteryV
}

func readAverageADC(adc machine.ADC, samples int) uint16 {
	var sum uint32 = 0
	for i := 0; i < samples; i++ {
		sum += uint32(adc.Get())
		time.Sleep(10 * time.Millisecond)
	}
	return uint16(sum / uint32(samples))
}
