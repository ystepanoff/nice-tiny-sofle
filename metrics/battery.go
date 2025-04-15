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
	nSamples             = 64
	adcResolution        = 16
	minBatteryV          = 3000
	maxBatteryV          = 4200
	batteryReadingPeriod = 10
)

func InitBattery() {
	machine.InitADC()
	batteryPin = machine.ADC{Pin: machine.P0_04}
	batteryPin.Configure(machine.ADCConfig{
		Reference:  minBatteryV,
		Resolution: adcResolution,
	})
	lastBatteryReading = time.Now()
	time.Sleep(1 * time.Second)
}

func ShouldReadBatteryLevel() bool {
	if time.Since(lastBatteryReading) >= batteryReadingPeriod*time.Second {
		lastBatteryReading = time.Now()
		return true
	}
	return false
}

func ReadBatteryLevel() float32 {
	return readBatteryVoltage()
}

func readBatteryVoltage() float32 {
	raw := readAverageADC(batteryPin, nSamples)
	//batteryV := float32(raw) * float32(3000.0/4096.0)
	/*if batteryV < minBatteryV {
		batteryV = minBatteryV
	} else if batteryV > maxBatteryV {
		batteryV = maxBatteryV
	}*/
	return float32(raw) / 2.0
}

func readAverageADC(adc machine.ADC, samples int) uint16 {
	var sum uint32 = 0
	for i := 0; i < samples; i++ {
		sum += uint32(adc.Get())
		time.Sleep(10 * time.Millisecond)
	}
	return uint16(sum / uint32(samples))
}
