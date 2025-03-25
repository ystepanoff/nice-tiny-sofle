package main

import (
	"fmt"
	"time"

	"github.com/ystepanoff/nice-tiny-sofle/display"
	"github.com/ystepanoff/nice-tiny-sofle/keyboard"
	"github.com/ystepanoff/nice-tiny-sofle/metrics"
	"github.com/ystepanoff/nice-tiny-sofle/sofle"
)

func updateDisplay() {
	display.Clear()
	if sofle.IsLeft {
		display.Draw(0, 0, 32, 32, sofle.SampleImage)
		display.Write(40, 16, fmt.Sprintf("%.2f", metrics.ReadBatteryLevel()))
	}
	display.Display()
}

func keyPressesLoop() {
	mat := keyboard.NewMatrix(sofle.RowPins, sofle.ColPins)
	for {
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

		time.Sleep(100 * time.Millisecond)
	}
}

func displayLoop() {
	display.Init()

	for {
		updateDisplay()

		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	// Battery readings are not reliable just yet,
	// leaving it for later to investigate
	// metrics.InitBattery()

	go keyPressesLoop()
	go displayLoop()
}
