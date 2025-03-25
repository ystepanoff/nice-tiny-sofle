package metrics

import (
	"machine"
	"time"
)

var batteryPin machine.ADC

const (
	nSamples      = 64
	adcResolution = 12
	minBatteryV   = 3000
	maxBatteryV   = 4200
)

func InitBattery() {
	machine.InitADC()
	batteryPin = machine.ADC{Pin: machine.P0_04}
	batteryPin.Configure(machine.ADCConfig{})
	batteryPin.Configure(machine.ADCConfig{})
	time.Sleep(1 * time.Second)
}

func ReadBatteryLevel() float32 {
	return readBatteryVoltage()
}

func readBatteryVoltage() float32 {
	time.Sleep(1 * time.Second)
	raw := readAverageADC(batteryPin, nSamples)
	batteryV := float32(raw) * float32(3300.0/4096.0) / float32(1.18)
	/*if batteryV < minBatteryV {
		batteryV = minBatteryV
	} else if batteryV > maxBatteryV {
		batteryV = maxBatteryV
	}*/
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
