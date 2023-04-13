package hal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hoani/chatbox/hal/ui"
)

type hal struct {
	button *button
	leds   *leds
	lcd    *lcd
}

func newHal() (*hal, error) {
	m := ui.NewUI()
	prog := tea.NewProgram(m)

	h := &hal{
		button: &button{},
		lcd:    &lcd{},
		leds:   &leds{prog: prog},
	}
	h.button.start()
	go prog.Run()
	return h, nil
}

func (h *hal) Leds() Leds {
	return h.leds
}

func (h *hal) LCD() LCD {
	return h.lcd
}

func (h *hal) Button() bool {
	return h.button.get()
}
