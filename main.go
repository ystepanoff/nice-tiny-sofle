package main

import (
	"fmt"
	"time"

	//"github.com/ystepanoff/nice-tiny-sofle/display"
	"github.com/ystepanoff/nice-tiny-sofle/keyboard"
	"github.com/ystepanoff/nice-tiny-sofle/metrics"
	"github.com/ystepanoff/nice-tiny-sofle/sofle"
)

/*
func updateDisplay() {
	display.Clear()
	if sofle.IsLeft {
		display.Draw(0, 0, 32, 32, sofle.SampleImage)
		level, err := metrics.ReadBatteryLevel()
		if err != nil {
			display.Write(40, 16, "ERR")
		} else {
			display.Write(40, 16, fmt.Sprintf("%.2f", float64(level)))
		}
	}
	display.Display()
}*/

func scanMatrix(mat *keyboard.Matrix) {
	pressed := mat.Scan()
	comb := make([][2]int, 0)
	for i := 0; i < len(pressed); i++ {
		for j := 0; j < len(pressed[i]); j++ {
			if pressed[i][j] {
				comb = append(comb, [2]int{i, j})
			}
		}
	}

	if len(comb) > 0 {
		fmt.Printf("%v\n\r", comb)
	}
}

func readBatteryLevel() {
	pct, err := metrics.ReadBatteryLevel()
	if err != nil {
		fmt.Printf("Battery reading error: %v\n\r", err)
		return
	}

	fmt.Printf("Battery: (%d%%)\n\r", pct)
}

func main() {
	metrics.InitSAADC()

	mat := keyboard.NewMatrix(sofle.RowPins, sofle.ColPins)

	// Create tickers for periodic tasks
	matrixTicker := time.NewTicker(100 * time.Millisecond)
	batteryTicker := time.NewTicker(metrics.BATTERY_READING_INTERVAL * time.Second)

	defer matrixTicker.Stop()
	defer batteryTicker.Stop()

	for {
		select {
		case <-matrixTicker.C:
			go scanMatrix(mat)
		case <-batteryTicker.C:
			go readBatteryLevel()
		}
	}
}
