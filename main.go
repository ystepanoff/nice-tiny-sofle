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
		display.Write(40, 16, fmt.Sprintf("%.2f", metrics.ReadBatteryLevel()))
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
	if metrics.ShouldReadBatteryLevel() {
		fmt.Printf("%v\n\r", metrics.ReadBatteryLevel())
	}
}

func main() {
	metrics.InitBattery()

	mat := keyboard.NewMatrix(sofle.RowPins, sofle.ColPins)

	for {
		go scanMatrix(mat)
		go readBatteryLevel()

		time.Sleep(100 * time.Millisecond)
	}
}
