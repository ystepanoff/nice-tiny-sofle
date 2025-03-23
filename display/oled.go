//go:build oled
// +build oled

package display

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

var (
	BLACK = color.RGBA{0, 0, 0, 255}
	WHITE = color.RGBA{255, 255, 255, 255}
)

var oled ssd1306.Device

func Init() {
	i2c := machine.I2C0

	machine.I2C0.Configure(machine.I2CConfig{
		SCL:       machine.P0_20,
		SDA:       machine.P0_17,
		Frequency: machine.TWI_FREQ_400KHZ,
	})

	oled = ssd1306.NewI2C(i2c)
	oled.Configure(ssd1306.Config{
		Width:   128,
		Height:  32,
		Address: 0x3C,
	})

	oled.ClearDisplay()
	oled.Display()
}

func Write(msg string, x, y int16) {
	tinyfont.WriteLine(&oled, &proggy.TinySZ8pt7b, x, y, msg, WHITE)
}

func Draw(x, y, w, h int16, data []byte) {
	byteWidth := (w + 7) / 8
	for py := int16(0); py < h; py++ {
		for px := int16(0); px < w; px++ {
			idx := int(py)*int(byteWidth) + int(px)/8
			bitMask := uint8(1 << (px & 7))
			if (data[idx] & bitMask) != 0 {
				oled.SetPixel(x+px, y+py, WHITE)
			}
		}
	}
	// oled.DrawMonochromeBitmap(x, y, w, h, data[:])
}
