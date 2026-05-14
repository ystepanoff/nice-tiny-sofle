package power

import (
	"sync"
	"sync/atomic"
	"time"

	"device/nrf"
	"machine"
)

type PowerState int32

const (
	Active PowerState = iota
	Idle
	Sleep
)

const (
	IDLE_TIMEOUT  = 30 * time.Second
	SLEEP_TIMEOUT = 5 * time.Minute
	WAKE_TIMEOUT  = 100 * time.Millisecond

	ACTIVE_SCAN_INTERVAL = 1 * time.Millisecond
	IDLE_SCAN_INTERVAL   = 10 * time.Millisecond
	SLEEP_SCAN_INTERVAL  = 100 * time.Millisecond

	MONITOR_POWER_STATE_INTERVAL = 100 * time.Millisecond

	ACTIVE_SPI_FREQ = 2000000
	IDLE_SPI_FREQ   = 1000000 // Sharp memory LCD min ~1 MHz
)

var (
	currentState     atomic.Int32 // PowerState
	lastActivityNano atomic.Int64
	onStateChange    func(PowerState)
	stopMonitor      chan struct{}
	transitionMu     sync.Mutex
)

func Init(onStateChangeCallback func(PowerState)) {
	onStateChange = onStateChangeCallback
	stopMonitor = make(chan struct{})
	currentState.Store(int32(Active))
	lastActivityNano.Store(time.Now().UnixNano())

	go monitorPowerState()
}

func monitorPowerState() {
	ticker := time.NewTicker(MONITOR_POWER_STATE_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			idleTime := time.Since(time.Unix(0, lastActivityNano.Load()))

			switch GetCurrentState() {
			case Active:
				if idleTime > IDLE_TIMEOUT {
					transitionTo(Idle)
				}
			case Idle:
				if idleTime > SLEEP_TIMEOUT {
					transitionTo(Sleep)
				} else if idleTime < WAKE_TIMEOUT {
					transitionTo(Active)
				}
			case Sleep:
				if idleTime < WAKE_TIMEOUT {
					transitionTo(Active)
				}
			}

		case <-stopMonitor:
			return
		}
	}
}

func transitionTo(newState PowerState) {
	transitionMu.Lock()
	defer transitionMu.Unlock()

	if PowerState(currentState.Load()) == newState {
		return
	}

	switch newState {
	case Active:
		nrf.CLOCK.EVENTS_HFCLKSTARTED.Set(0)
		nrf.CLOCK.TASKS_HFCLKSTART.Set(1)
		for nrf.CLOCK.EVENTS_HFCLKSTARTED.Get() == 0 {
		}

		machine.SPI0.Configure(machine.SPIConfig{
			Frequency: ACTIVE_SPI_FREQ,
			SCK:       machine.P0_20,
			SDO:       machine.P0_17,
			SDI:       machine.P0_25,
			Mode:      0,
			LSBFirst:  true,
		})
	case Idle:
		machine.SPI0.Configure(machine.SPIConfig{
			Frequency: IDLE_SPI_FREQ,
			SCK:       machine.P0_20,
			SDO:       machine.P0_17,
			SDI:       machine.P0_25,
			Mode:      0,
			LSBFirst:  true,
		})
	case Sleep:
		nrf.CLOCK.TASKS_HFCLKSTOP.Set(1)
		configureWakeSources()
	}

	currentState.Store(int32(newState))
	if onStateChange != nil {
		onStateChange(newState)
	}
}

func UpdateActivity() {
	lastActivityNano.Store(time.Now().UnixNano())
	if GetCurrentState() != Active {
		transitionTo(Active)
	}
}

func GetCurrentState() PowerState {
	return PowerState(currentState.Load())
}

func GetScanInterval() time.Duration {
	switch GetCurrentState() {
	case Idle:
		return IDLE_SCAN_INTERVAL
	case Sleep:
		return SLEEP_SCAN_INTERVAL
	default:
		return ACTIVE_SCAN_INTERVAL
	}
}

func configureWakeSources() {
	// placeholder
}

func Cleanup() {
	if stopMonitor != nil {
		close(stopMonitor)
	}
}
