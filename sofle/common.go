package sofle

import (
	"machine"

	kb "machine/usb/hid/keyboard"
)

var (
	RowPins = []machine.Pin{
		machine.P0_24,
		machine.P1_00,
		machine.P0_11,
		machine.P1_04,
		machine.P1_06,
	}

	ColPins = []machine.Pin{
		machine.P0_02,
		machine.P1_15,
		machine.P1_13,
		machine.P1_11,
		machine.P0_10,
		machine.P0_09,
	}

	KeyNo = kb.Keycode(0x00)

	Keymap [][]kb.Keycode
)
