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

func (l *lcd) Write(line1, line2 string, color *RGB) {
	var lcdColor *lipgloss.Color
	if color != nil {
		c := RGBtoLipglossColor(*color)
		lcdColor = &c
	}
	l.prog.Send(ui.LCDUpdate{
		Color: lcdColor,
		Lines: [2]string{
			line1 + " ",
			line2 + " ",
		},
	})
}
