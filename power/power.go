package power

import (
	"device/nrf"
	"machine"
	"runtime"
	"time"
)

type PowerState int

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
	IDLE_SPI_FREQ   = 500000
)

var (
	currentState     PowerState = Active
	lastActivity     time.Time  = time.Now()
	onStateChange    func(PowerState)
	stopMonitor      chan struct{}
	lastScanInterval time.Duration
)

func Init(onStateChangeCallback func(PowerState)) {
	onStateChange = onStateChangeCallback
	stopMonitor = make(chan struct{})
	lastScanInterval = ACTIVE_SCAN_INTERVAL

	configurePowerManagement()
	go monitorPowerState()
}

func configurePowerManagement() {
	nrf.POWER.SYSTEMOFF.Set(0)
	nrf.POWER.POFCON.Set(0)
	nrf.POWER.RESETREAS.Set(0)
}

func monitorPowerState() {
	ticker := time.NewTicker(MONITOR_POWER_STATE_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			idleTime := now.Sub(lastActivity)

			switch currentState {
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

			runtime.GC()

		case <-stopMonitor:
			return
		}
	}
}

func transitionTo(newState PowerState) {
	if newState == currentState {
		return
	}

	switch currentState {
	case Active:
	case Idle:
	case Sleep:
	}

	switch newState {
	case Active:
		nrf.CLOCK.EVENTS_HFCLKSTARTED.Set(0)
		nrf.CLOCK.TASKS_HFCLKSTART.Set(1)

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

	currentState = newState
	if onStateChange != nil {
		onStateChange(newState)
	}
}

func UpdateActivity() {
	lastActivity = time.Now()
	if currentState != Active {
		transitionTo(Active)
	}
}

func GetCurrentState() PowerState {
	return currentState
}

func GetScanInterval() time.Duration {
	interval := ACTIVE_SCAN_INTERVAL
	switch currentState {
	case Active:
		interval = ACTIVE_SCAN_INTERVAL
	case Idle:
		interval = IDLE_SCAN_INTERVAL
	case Sleep:
		interval = SLEEP_SCAN_INTERVAL
	}

	if interval != lastScanInterval {
		lastScanInterval = interval
	}

	return interval
}

func configureWakeSources() {
	// placeholder
}

func Cleanup() {
	if stopMonitor != nil {
		close(stopMonitor)
	}
}
