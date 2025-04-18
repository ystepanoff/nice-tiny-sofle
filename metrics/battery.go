package metrics

import (
	"machine"
	"time"
)

var (
	batteryPin         machine.ADC
	lastBatteryReading time.Time
)

const (
	adcReference         = 600
	adcResolution        = 12
	adcSamples           = 16
	adcSampleTime        = 40
	minBatteryV          = 3000
	maxBatteryV          = 4200
	batteryReadingPeriod = 60
)

func InitBattery() {
	machine.InitADC()
	batteryPin = machine.ADC{Pin: machine.AIN2}
	batteryPin.Configure(machine.ADCConfig{
		Reference:  adcReference,
		Resolution: adcReference,
		Samples:    adcSamples,
		SampleTime: adcSampleTime,
	})
	time.Sleep(1 * time.Second)
	lastBatteryReading = time.Now()
}

func ShouldReadBatteryLevel() bool {
	if time.Since(lastBatteryReading) >= batteryReadingPeriod*time.Second {
		lastBatteryReading = time.Now()
		return true
	}
	return false
}

func ReadBatteryLevel() uint16 {
	return readBatteryVoltage()
}

func readBatteryVoltage() uint16 {
	time.Sleep(1 * time.Second)
	raw := batteryPin.Get()
	//raw := readAverageADC(batteryPin, nSamples)
	//batteryV := float32(raw) * float32(3000.0/4096.0)
	/*if batteryV < minBatteryV {
		batteryV = minBatteryV
	} else if batteryV > maxBatteryV {
		batteryV = maxBatteryV
	}*/
	return raw
}
