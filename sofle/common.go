package sofle

import "github.com/tinygo-org/tinygo/src/machine"

type Layer int

const (
	Base Layer = iota
	Upper
	Lower
	Adjust
)

var (
	RowPins = []machine.Pin{
		machine.P0_24,
		machine.P1_00,
		machine.P0_11,
		machine.P1_04,
		machine.P1_06,
	}

	RolPins = []machine.Pin{
		machine.P0_02,
		machine.P1_15,
		machine.P1_13,
		machine.P1_11,
		machine.P0_10,
		machine.P0_09,
	}
)
