package metrics

import (
	"device/nrf"
	"fmt"
	"time"
	"unsafe"
)

const (
	// SAADC oversample register value. Hardware does (1 << OVERSAMPLE_REG) samples
	// per RESULT word — here, 16× oversampling.
	OVERSAMPLE_REG    = 4
	OVERSAMPLE_FACTOR = 1 << OVERSAMPLE_REG

	VDDHDIV5                 = 0x0D // SAADC_CH_PSELP_PSELP_VDDHDIV5
	BATTERY_READING_INTERVAL = 60 * time.Second

	SAADC_CALIBRATION_TIMEOUT = 100 * time.Millisecond
	SAADC_EVENT_TIMEOUT       = 10 * time.Millisecond
	SAADC_RESET_DELAY         = 1 * time.Millisecond
)

var (
	saadcWord   uint32
	bufPtr      uint32
	lastReading time.Time
	lastLevel   uint16
)

func InitSAADC() error {
	nrf.SAADC.ENABLE.Set(1)
	nrf.SAADC.RESOLUTION.Set(nrf.SAADC_RESOLUTION_VAL_12bit)
	nrf.SAADC.OVERSAMPLE.Set(OVERSAMPLE_REG)

	nrf.SAADC.CH[0].PSELP.Set(VDDHDIV5)
	nrf.SAADC.CH[0].PSELN.Set(0x1F)

	cfg := uint32(
		nrf.SAADC_CH_CONFIG_GAIN_Gain1_2<<nrf.SAADC_CH_CONFIG_GAIN_Pos |
			nrf.SAADC_CH_CONFIG_REFSEL_Internal<<nrf.SAADC_CH_CONFIG_REFSEL_Pos |
			nrf.SAADC_CH_CONFIG_TACQ_40us<<nrf.SAADC_CH_CONFIG_TACQ_Pos |
			nrf.SAADC_CH_CONFIG_MODE_SE<<nrf.SAADC_CH_CONFIG_MODE_Pos,
	)
	nrf.SAADC.CH[0].CONFIG.Set(cfg)

	nrf.SAADC.EVENTS_CALIBRATEDONE.Set(0)
	nrf.SAADC.TASKS_CALIBRATEOFFSET.Set(1)

	deadline := time.Now().Add(SAADC_CALIBRATION_TIMEOUT)
	for nrf.SAADC.EVENTS_CALIBRATEDONE.Get() == 0 {
		if time.Now().After(deadline) {
			return fmt.Errorf("SAADC calibration timeout")
		}
	}

	bufPtr = uint32(uintptr(unsafe.Pointer(&saadcWord)))
	return nil
}

func resetSAADC() {
	nrf.SAADC.ENABLE.Set(0)
	time.Sleep(SAADC_RESET_DELAY)
	nrf.SAADC.ENABLE.Set(1)
	time.Sleep(SAADC_RESET_DELAY)
}

func waitForEvent(reg interface{ Get() uint32 }, name string) error {
	deadline := time.Now().Add(SAADC_EVENT_TIMEOUT)
	for reg.Get() == 0 {
		if time.Now().After(deadline) {
			resetSAADC()
			return fmt.Errorf("timeout waiting for SAADC %s event", name)
		}
	}
	return nil
}

func sampleRaw() (uint16, error) {
	nrf.SAADC.RESULT.PTR.Set(bufPtr)
	nrf.SAADC.RESULT.MAXCNT.Set(1)

	nrf.SAADC.EVENTS_STARTED.Set(0)
	nrf.SAADC.EVENTS_END.Set(0)
	nrf.SAADC.EVENTS_STOPPED.Set(0)

	nrf.SAADC.TASKS_START.Set(1)
	if err := waitForEvent(&nrf.SAADC.EVENTS_STARTED, "STARTED"); err != nil {
		return 0, err
	}

	// One TASKS_SAMPLE per oversample step. The hardware accumulates internally
	// and raises EVENTS_END after OVERSAMPLE_FACTOR samples.
	for i := 0; i < OVERSAMPLE_FACTOR; i++ {
		nrf.SAADC.TASKS_SAMPLE.Set(1)
	}

	if err := waitForEvent(&nrf.SAADC.EVENTS_END, "END"); err != nil {
		return 0, err
	}

	nrf.SAADC.TASKS_STOP.Set(1)
	if err := waitForEvent(&nrf.SAADC.EVENTS_STOPPED, "STOPPED"); err != nil {
		return 0, err
	}

	return uint16(saadcWord & 0xFFFF), nil
}

func rawToMillivolts(r uint16) uint32 {
	return uint32(r) * 6000 >> 12
}

// Same transformation as used in ZMK
func liIonPct(mV uint32) uint8 {
	switch {
	case mV >= 4200:
		return 100
	case mV <= 3450:
		return 0
	default:
		return uint8(int(mV)*2/15 - 459)
	}
}

func ReadBatteryLevel() (uint16, error) {
	if !lastReading.IsZero() && time.Since(lastReading) < BATTERY_READING_INTERVAL/2 {
		return lastLevel, nil
	}

	raw, err := sampleRaw()
	if err != nil {
		return lastLevel, err
	}

	level := uint16(liIonPct(rawToMillivolts(raw)))
	lastLevel = level
	lastReading = time.Now()
	return level, nil
}
