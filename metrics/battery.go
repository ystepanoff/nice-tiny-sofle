package metrics

import (
	"device/nrf"
	"fmt"
	"time"
	"unsafe"
)

const (
	OVERSAMPLE               = 4    // 4x oversampling
	VDDHDIV5                 = 0x0D // SAADC_CH_PSELP_PSELP_VDDHDIV5
	BATTERY_READING_INTERVAL = 5
	SAADC_TIMEOUT            = 2 * time.Second
	SAADC_SAMPLE_DELAY       = 10 * time.Millisecond
	SAADC_RESET_DELAY        = 10 * time.Millisecond
)

var (
	saadcWord uint32
	bufPtr    uint32
)

func InitSAADC() {
	nrf.SAADC.ENABLE.Set(1)
	nrf.SAADC.RESOLUTION.Set(nrf.SAADC_RESOLUTION_VAL_12bit)
	nrf.SAADC.OVERSAMPLE.Set(OVERSAMPLE)

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
	for nrf.SAADC.EVENTS_CALIBRATEDONE.Get() == 0 {
	}

	bufPtr = uint32(uintptr(unsafe.Pointer(&saadcWord)))
}

func resetSAADC() {
	nrf.SAADC.ENABLE.Set(0)
	time.Sleep(SAADC_RESET_DELAY)
	nrf.SAADC.ENABLE.Set(1)
	time.Sleep(SAADC_RESET_DELAY)
}

func sampleRaw() (uint16, error) {
	nrf.SAADC.RESULT.PTR.Set(bufPtr)
	nrf.SAADC.RESULT.MAXCNT.Set(1)

	nrf.SAADC.EVENTS_STARTED.Set(0)
	nrf.SAADC.EVENTS_END.Set(0)
	nrf.SAADC.EVENTS_STOPPED.Set(0)

	nrf.SAADC.TASKS_START.Set(1)

	startTime := time.Now()
	for nrf.SAADC.EVENTS_STARTED.Get() == 0 {
		if time.Since(startTime) > SAADC_TIMEOUT {
			resetSAADC()
			return 0, fmt.Errorf("timeout waiting for SAADC STARTED event")
		}
	}

	for i := 0; i < (1 << OVERSAMPLE); i++ {
		nrf.SAADC.TASKS_SAMPLE.Set(1)
		time.Sleep(SAADC_SAMPLE_DELAY)
	}

	startTime = time.Now()
	for nrf.SAADC.EVENTS_END.Get() == 0 {
		if time.Since(startTime) > SAADC_TIMEOUT {
			resetSAADC()
			return 0, fmt.Errorf("timeout waiting for SAADC END event")
		}
		time.Sleep(10 * time.Millisecond)
	}

	nrf.SAADC.TASKS_STOP.Set(1)

	startTime = time.Now()
	for nrf.SAADC.EVENTS_STOPPED.Get() == 0 {
		if time.Since(startTime) > SAADC_TIMEOUT {
			resetSAADC()
			return 0, fmt.Errorf("timeout waiting for SAADC STOPPED event")
		}
		time.Sleep(10 * time.Millisecond)
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
	raw, err := sampleRaw()
	if err != nil {
		return 0, err
	}
	return uint16(liIonPct(rawToMillivolts(raw))), nil
}
