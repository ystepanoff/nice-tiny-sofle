package main

import (
	"time"

	"github.com/ystepanoff/nice-tiny-sofle/display"
	"github.com/ystepanoff/nice-tiny-sofle/keyboard"
	"github.com/ystepanoff/nice-tiny-sofle/metrics"
	"github.com/ystepanoff/nice-tiny-sofle/sofle"
)

func main() {
	if err := display.Init(); err != nil {
		println("Failed to initialise display:", err.Error())
		return
	}
	println("nice!view initialised")

	metrics.InitSAADC()
	println("SAADC initialised")

	// Initialize keyboard matrix and handler
	matrix := keyboard.NewMatrix(sofle.RowPins, sofle.ColPins)
	keyHandler := keyboard.NewKeyHandler(matrix, func(row, col int, pressed bool) {
		println("Key event:", row, col, pressed)
	})
	println("Keyboard initialized")

	batteryTicker := time.NewTicker(metrics.BATTERY_READING_INTERVAL)
	defer batteryTicker.Stop()

	keyboardTicker := time.NewTicker(keyboard.SCAN_INTERVAL)
	defer keyboardTicker.Stop()

	// Main loop
	for {
		select {
		case <-batteryTicker.C:
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

		case <-keyboardTicker.C:
			keyHandler.Scan()
		}
	}
}
