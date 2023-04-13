package hal

import (
	"fmt"
	"math"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hoani/chatbox/hal/ui"
)

type leds struct {
	value bool
	lock  sync.Mutex
	prog  *tea.Program
}

func newLeds() *leds {
	return &leds{}
}

func (l *leds) HSV(index int, values ...HSV) {
	rgbs := make([]RGB, len(values))
	for i, value := range values {
		rgbs[i] = HSVtoRGB(value)
	}
	l.RGB(index, rgbs...)
}

func (l *leds) RGB(index int, values ...RGB) {

	colors := make([]lipgloss.Color, len(values))
	for i, value := range values {
		colors[i] = lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", value.R, value.G, value.B))
	}

	l.prog.Send(ui.LEDColors{
		Index:  index,
		Colors: colors,
	})

}
func (l *leds) Show() {
	l.prog.Send(ui.LEDShow{})
}
func (l *leds) Clear() {
	l.prog.Send(ui.LEDClear{})

}

// https://www.rapidtables.com/convert/color/hsv-to-rgb.html
func HSVtoRGB(in HSV) RGB {
	c := float64(in.V) / 256.0 * float64(in.S) / 256.0
	theta := 360 * float64(in.H) / 256
	x := c * math.Abs(float64(int(theta/60.0)%2)-1)
	m := float64(in.V) - c

	c = c * 256
	x = x * 256
	m = m * 256

	if theta < 60 {
		return RGB{R: uint8(c + m), G: uint8(x + m), B: uint8(m)}
	} else if theta < 120 {
		return RGB{R: uint8(x + m), G: uint8(c + m), B: uint8(m)}
	} else if theta < 180 {
		return RGB{R: uint8(m), G: uint8(c + m), B: uint8(x + m)}
	} else if theta < 240 {
		return RGB{R: uint8(m), G: uint8(x + m), B: uint8(c + m)}
	} else if theta < 300 {
		return RGB{R: uint8(x + m), G: uint8(m), B: uint8(c + m)}
	} else {
		return RGB{R: uint8(c + m), G: uint8(m), B: uint8(x + m)}
	}
}
