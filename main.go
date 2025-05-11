//go:build niceview && (left || right)
// +build niceview
// +build left right

package main

import (
	"time"

	"machine/usb/hid/keyboard"

	"github.com/ystepanoff/nice-tiny-sofle/display"
	kb "github.com/ystepanoff/nice-tiny-sofle/keyboard"
	"github.com/ystepanoff/nice-tiny-sofle/metrics"
	"github.com/ystepanoff/nice-tiny-sofle/power"
	"github.com/ystepanoff/nice-tiny-sofle/sofle"
)

const (
	INITIAL_SLEEP_INTERVAL               = 3 * time.Second
	SCAN_INTERVAL_MONITOR_SLEEP_INTERVAL = 100 * time.Millisecond
)

var hidkb = keyboard.New()

func main() {
	time.Sleep(INITIAL_SLEEP_INTERVAL)

	metrics.InitSAADC()
	println("SAADC initialised")

	if err := display.Init(); err != nil {
		println("Failed to initialise display:", err.Error())
		return
	}
	println("nice!view initialised")

	power.Init(func(state power.PowerState) {
		println("Power state changed to:", state)
	})
	println("Power management initialised")

	matrix := kb.NewMatrix(sofle.RowPins, sofle.ColPins)
	keyHandler := kb.NewKeyHandler(matrix, func(row, col int, pressed bool) {
		println("Key event:", row, col, pressed)
		power.UpdateActivity()

		if row < len(sofle.Keymap) && col < len(sofle.Keymap[row]) {
			key := sofle.Keymap[row][col]
			if pressed {
				hidkb.Down(key)
			} else {
				hidkb.Up(key)
			}
		}
	})
	println("Keyboard initialised")

	// Battery reading goroutine
	go func() {
		ticker := time.NewTicker(metrics.BATTERY_READING_INTERVAL)
		defer ticker.Stop()

		for {
			<-ticker.C
			if power.GetCurrentState() != power.Sleep {
				level, err := metrics.ReadBatteryLevel()
				if err != nil {
					println("Failed to read battery level:", err.Error())
					continue
				}
				println("Battery level:", level, "%")

				if err := display.Update(uint8(level)); err != nil {
					println("Failed to update display:", err.Error())
					continue
				}
			}
		}
	}()

	keyboardTicker := time.NewTicker(power.GetScanInterval())
	defer keyboardTicker.Stop()

	// Scan interval monitor
	go func() {
		var scanInterval time.Duration
		for {
			time.Sleep(SCAN_INTERVAL_MONITOR_SLEEP_INTERVAL)
			newInterval := power.GetScanInterval()
			if newInterval != scanInterval {
				keyboardTicker.Reset(newInterval)
				println("Keyboard scan interval updated to:", newInterval)
				scanInterval = newInterval
			}
		}
	}()

	// Main loop - only handles keyboard scanning
	for {
		<-keyboardTicker.C
		keyHandler.Scan()
	}
}
