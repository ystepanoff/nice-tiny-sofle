package keyboard

import (
	"time"
)

const (
	SCAN_INTERVAL   = 100 * time.Microsecond
	REPEAT_DELAY    = 500 * time.Millisecond // Time before repeat starts
	REPEAT_INTERVAL = 50 * time.Millisecond  // Time between repeats
)

type KeyState struct {
	Pressed     bool
	LastPressed bool
	RepeatCount int
	LastRepeat  time.Time
}

type KeyHandler struct {
	matrix     *Matrix
	keyStates  [][]KeyState
	lastScan   time.Time
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
		lastScan:   time.Now(),
		onKeyEvent: onKeyEvent,
	}
}

func (h *KeyHandler) Scan() {
	now := time.Now()
	if now.Sub(h.lastScan) < SCAN_INTERVAL {
		return
	}
	h.lastScan = now

	pressed := h.matrix.Scan()

	for r := range pressed {
		for c := range pressed[r] {
			state := &h.keyStates[r][c]
			current := pressed[r][c]

			// State change detection
			if current != state.Pressed {
				state.Pressed = current
				if state.Pressed != state.LastPressed {
					state.LastPressed = state.Pressed
					if h.onKeyEvent != nil {
						h.onKeyEvent(r, c, state.Pressed)
					}
					if state.Pressed {
						state.LastRepeat = now
					} else {
						state.RepeatCount = 0
					}
				}
			}

			// Key repeat handling
			if state.Pressed {
				timeSinceLastRepeat := now.Sub(state.LastRepeat)

				if state.RepeatCount == 0 {
					if timeSinceLastRepeat >= REPEAT_DELAY {
						state.RepeatCount++
						state.LastRepeat = now
						if h.onKeyEvent != nil {
							h.onKeyEvent(r, c, true)
						}
					}
				} else if timeSinceLastRepeat >= REPEAT_INTERVAL {
					state.LastRepeat = now
					if h.onKeyEvent != nil {
						h.onKeyEvent(r, c, true)
					}
				}
			}
		}
	}
}
