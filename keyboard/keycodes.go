package keyboard

type Keycode uint16

// Basic keyboard usage IDs from 0x04 through 0x73:
const (
	KC_NO Keycode = 0x00
	KC_A  Keycode = 0x04
	KC_B  Keycode = 0x05
	KC_C  Keycode = 0x06
	KC_D  Keycode = 0x07
	KC_E  Keycode = 0x08
	KC_F  Keycode = 0x09
	KC_G  Keycode = 0x0A
	KC_H  Keycode = 0x0B
	KC_I  Keycode = 0x0C
	KC_J  Keycode = 0x0D
	KC_K  Keycode = 0x0E
	KC_L  Keycode = 0x0F
	KC_M  Keycode = 0x10
	KC_N  Keycode = 0x11
	KC_O  Keycode = 0x12
	KC_P  Keycode = 0x13
	KC_Q  Keycode = 0x14
	KC_R  Keycode = 0x15
	KC_S  Keycode = 0x16
	KC_T  Keycode = 0x17
	KC_U  Keycode = 0x18
	KC_V  Keycode = 0x19
	KC_W  Keycode = 0x1A
	KC_X  Keycode = 0x1B
	KC_Y  Keycode = 0x1C
	KC_Z  Keycode = 0x1D

	KC_1 Keycode = 0x1E
	KC_2 Keycode = 0x1F
	KC_3 Keycode = 0x20
	KC_4 Keycode = 0x21
	KC_5 Keycode = 0x22
	KC_6 Keycode = 0x23
	KC_7 Keycode = 0x24
	KC_8 Keycode = 0x25
	KC_9 Keycode = 0x26
	KC_0 Keycode = 0x27

	KC_ENTER  Keycode = 0x28
	KC_ESC    Keycode = 0x29
	KC_BSPACE Keycode = 0x2A
	KC_TAB    Keycode = 0x2B
	KC_SPACE  Keycode = 0x2C
	// ...
	KC_F1 Keycode = 0x3A
	KC_F2 Keycode = 0x3B
	// ...
	KC_F12 Keycode = 0x45
	// ...
	KC_PGUP  Keycode = 0x4B
	KC_PGDN  Keycode = 0x4E
	KC_RIGHT Keycode = 0x4F
	KC_LEFT  Keycode = 0x50
	KC_DOWN  Keycode = 0x51
	KC_UP    Keycode = 0x52
	// ...
	KC_F24 Keycode = 0x73
	// etc.
)

// Modifiers
const (
	KC_LCTRL  Keycode = 0xE0
	KC_LSHIFT Keycode = 0xE1
	KC_LALT   Keycode = 0xE2
	KC_LGUI   Keycode = 0xE3
	KC_RCTRL  Keycode = 0xE4
	KC_RSHIFT Keycode = 0xE5
	KC_RALT   Keycode = 0xE6
	KC_RGUI   Keycode = 0xE7
)

// Media keys
const (
	KC_VOLU  Keycode = 0x0100 + 0xE9
	KC_VOLD  Keycode = 0x0100 + 0xEA
	KC_MUTE  Keycode = 0x0100 + 0xE2
	KC_MNEXT Keycode = 0x0100 + 0xB5
	KC_MPREV Keycode = 0x0100 + 0xB6
	KC_MSTOP Keycode = 0x0100 + 0xB7
	// ...
)
