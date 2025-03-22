//go:build left
// +build left

package sofle

import (
	kb "github.com/ystepanoff/nice-tiny-sofle/keyboard"
)

const IsLeft = true

var keymaps = map[Layer][][]kb.Keycode{
	Base: {
		{kb.KC_Q, kb.KC_W, kb.KC_E, kb.KC_R, kb.KC_T},
		{kb.KC_A, kb.KC_S, kb.KC_D, kb.KC_F, kb.KC_G},
		{kb.KC_Z, kb.KC_X, kb.KC_C, kb.KC_V, kb.KC_B},
	},
	Upper: {
		{kb.KC_1, kb.KC_2, kb.KC_3, kb.KC_4, kb.KC_5},
		{kb.KC_6, kb.KC_7, kb.KC_8, kb.KC_9, kb.KC_0},
		{kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO},
	},
	Lower: {
		{kb.KC_LEFT, kb.KC_UP, kb.KC_RIGHT, kb.KC_DOWN, kb.KC_ESC},
		{kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO},
		{kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO},
	},
	Adjust: {
		{kb.KC_VOLU, kb.KC_VOLD, kb.KC_PGUP, kb.KC_PGDN, kb.KC_ESC},
		{kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO},
		{kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO, kb.KC_NO},
	},
}
