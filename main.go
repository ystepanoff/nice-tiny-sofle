package main

import (
	"fmt"
	"time"

	"github.com/ystepanoff/nice-tiny-sofle/display"
	"github.com/ystepanoff/nice-tiny-sofle/keyboard"
	"github.com/ystepanoff/nice-tiny-sofle/sofle"
)

func updateDisplay() {
	if sofle.IsLeft {
		display.Draw(0, 0, 32, 32, sofle.SampleImage)
	}
}

func main() {
	display.Init()
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

		updateDisplay()

		time.Sleep(100 * time.Millisecond)

		if len(comb) > 0 {
			fmt.Printf("%v\n\r", comb)
		}
	}
}
