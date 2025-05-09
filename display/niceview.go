package display

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/sharpmem"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

const (
	DISPLAY_INIT_DELAY = 100 * time.Millisecond
	DISPLAY_WIDTH      = 160
	DISPLAY_HEIGHT     = 68
	DISPLAY_FREQUENCY  = 2000000
)

type DisplayState struct {
	BatteryLevel uint8
}

var displayState DisplayState
var display sharpmem.Device

func hasDisplayStateChanged(newState DisplayState) bool {
	return displayState.BatteryLevel != newState.BatteryLevel
}

func Init() error {
	machine.P0_06.Configure(machine.PinConfig{Mode: machine.PinOutput})

	err := machine.SPI0.Configure(machine.SPIConfig{
		Frequency: DISPLAY_FREQUENCY,
		SCK:       machine.P0_20,
		SDO:       machine.P0_17,
		SDI:       machine.P0_25,
		Mode:      0,
		LSBFirst:  true,
	})
	if err != nil {
		return err
	}

	display = sharpmem.New(machine.SPI0, machine.P0_06)

	display.Configure(sharpmem.Config{
		Width:  DISPLAY_WIDTH,
		Height: DISPLAY_HEIGHT,
	})

	display.Clear()
	display.Display()
	time.Sleep(DISPLAY_INIT_DELAY)

	return nil
}

func Clear() error {
	return display.Clear()
}

func Display() error {
	return display.Display()
}

func Update(batteryLevel uint8) error {
	newDisplayState := DisplayState{
		BatteryLevel: batteryLevel,
	}
	if hasDisplayStateChanged(newDisplayState) {
		displayState = newDisplayState

		if err := Clear(); err != nil {
			return err
		}
		text := fmt.Sprintf("%d%%", batteryLevel)
		tinyfont.WriteLine(&display, &proggy.TinySZ8pt7b, 40, 40, text, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		return Display()
	}
	return nil
}
