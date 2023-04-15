package hal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hoani/chatbox/hal/ui"
)

type lcd struct {
	prog *tea.Program
}

func newLCD(prog *tea.Program) *lcd {
	return &lcd{prog: prog}
}

func (l *lcd) Write(line1, line2 string, color LCDColor) {
	c := LCDtoLipglossColor(color)
	l.prog.Send(ui.LCDUpdate{
		Color: &c,
		Lines: [2]string{
			line1 + " ",
			line2 + " ",
		},
	})
}

func LCDtoLipglossColor(in LCDColor) lipgloss.Color {
	switch in {
	case LCDRed:
		return lipgloss.Color("#ff0000")
	case LCDGreen:
		return lipgloss.Color("#00ff00")
	case LCDAqua:
		return lipgloss.Color("#00ffff")
	case LCDBlue:
		return lipgloss.Color("#0000ff")
	default:
		return lipgloss.Color("#444444")
	}
}
