package keyboard

import (
	"time"

	"github.com/tinygo-org/tinygo/src/machine"
)

type Matrix struct {
	rowPins []machine.Pin
	colPins []machine.Pin
}

func NewMatrix(rowPins, colPins []machine.Pin) *Matrix {
	for _, rp := range rowPins {
		rp.Configure(machine.PinConfig{Mode: machine.PinOutput})
		rp.High()
	}
	for _, cp := range colPins {
		cp.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	return &Matrix{
		rowPins: rowPins,
		colPins: colPins,
	}
}

func (m *Matrix) Scan() [][]bool {
	pressed := make([][]bool, len(m.rowPins))

	for r, row := range m.rowPins {
		pressed[r] = make([]bool, len(m.colPins))

		row.Low()
		time.Sleep(*time.Microsecond)

		for c, col := range m.colPins {
			pressed[c] = !col.Get()
		}

		row.High()
	}

	return pressed
}
