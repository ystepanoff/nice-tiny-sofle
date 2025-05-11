//go:build right
// +build right

package sofle

const IsLeft = false

func init() {
	Keymap = [][]kb.Keycode{
		{kb.Key6, kb.Key7, kb.Key8, kb.Key9, kb.Key0, kb.KeyEscape},
		{kb.KeyJ, kb.KeyL, kb.KeyU, kb.KeyY, kb.KeySemicolon, kb.KeyBackslash},
		{kb.KeyM, kb.KeyN, kb.KeyE, kb.KeyI, kb.KeyO, kb.KeyQuote},
		{kb.KeyK, kb.KeyH, kb.KeyComma, kb.KeyPeriod, kb.KeySlash, kb.KeyRightShift, KeyNo},
		{kb.KeySpace, KeyNo, kb.KeyModifierRightGUI, kb.KeyModifierRightAlt, kb.KeyModifierRightCtrl},
	}
}
