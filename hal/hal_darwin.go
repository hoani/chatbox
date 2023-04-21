package hal

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hoani/chatbox/hal/ui"
)

type hal struct {
	button *button
	leds   *leds
	lcd    *lcd
	prog   *tea.Program
	m      *ui.Model
}

func newHal() (*hal, error) {
	m := ui.NewUI()
	prog := tea.NewProgram(m)

	h := &hal{
		lcd:  newLCD(prog),
		leds: newLeds(prog),
		prog: prog,
		m:    m,
	}

	go prog.Run()
	return h, nil
}

func (h *hal) Debug(msg string) {
	h.prog.Send(ui.Debug(msg))
}

func (h *hal) Leds() Leds {
	return h.leds
}

func (h *hal) LCD() LCD {
	return h.lcd
}

func (h *hal) Button() bool {
	return h.m.ButtonState()
}

func (h *hal) Shutdown() {
	os.Exit(0)
}
