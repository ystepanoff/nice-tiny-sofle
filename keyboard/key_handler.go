package keyboard

type KeyState struct {
	Pressed bool
}

type KeyHandler struct {
	matrix     *Matrix
	keyStates  [][]KeyState
	onKeyEvent func(row, col int, pressed bool)
}

func NewKeyHandler(matrix *Matrix, onKeyEvent func(row, col int, pressed bool)) *KeyHandler {
	rows := len(matrix.rowPins)
	cols := len(matrix.colPins)

	keyStates := make([][]KeyState, rows)
	for i := range keyStates {
		keyStates[i] = make([]KeyState, cols)
	}

	return &KeyHandler{
		matrix:     matrix,
		keyStates:  keyStates,
		onKeyEvent: onKeyEvent,
	}
}

func (h *KeyHandler) Scan() {
	pressed := h.matrix.Scan()

	for r := range pressed {
		for c := range pressed[r] {
			state := &h.keyStates[r][c]
			current := pressed[r][c]

			if current != state.Pressed {
				state.Pressed = current
				if h.onKeyEvent != nil {
					h.onKeyEvent(r, c, current)
				}
			}
		}
	}
}
