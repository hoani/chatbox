package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

var viewCount = 0
var updateCount = 0

const (
	padding   = 2
	maxWidth  = 80
	nLEDs     = 24
	LCDWidth  = 16
	LCDHeight = 2
)

type LEDColors struct {
	Index  int
	Colors []lipgloss.Color
}

type LEDShow struct{}

type LEDClear struct{}

type LCDColor lipgloss.Color
type LCDLine1 string
type LCDLine2 string

type model struct {
	pending  [nLEDs]lipgloss.Color
	leds     [nLEDs]string
	lcdColor lipgloss.Color
	lcdLine1 string
	lcdLine2 string
}

func NewUI() tea.Model {
	colors := [nLEDs]lipgloss.Color{}
	leds := [nLEDs]string{}
	for i := range colors {
		colors[i] = lipgloss.Color("#000000")
		leds[i] = lipgloss.NewStyle().Background(colors[i]).Render(" ")
	}

	return model{
		pending:  colors,
		leds:     leds,
		lcdColor: lipgloss.Color("#0000ee"),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateCount++
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case LEDColors:
		for i, color := range msg.Colors {
			if i+msg.Index < nLEDs {
				m.pending[i] = color
			}
		}
		return m, nil

	case LEDShow:
		for i, color := range m.pending {
			m.leds[i] = lipgloss.NewStyle().Background(color).Render(" ")
		}
		return m, nil

	case LEDClear:
		for i := range m.pending {
			m.pending[i] = lipgloss.Color("#000000")
		}
		return m, nil

	}
	return m, nil
}

func (m model) View() string {

	return "\n" +
		m.ViewLEDs() + "\n\n" +
		helpStyle("Press any key to quit")
}

func (m model) ViewLEDs() string {
	return fmt.Sprintf(`
   %s%s%s%s
  %s    %s
 %s      %s
%s        %s
%s        %s
%s        %s
%s        %s
 %s      %s
  %s    %s
   %s%s%s%s
`, m.leds[0], m.leds[1], m.leds[2], m.leds[3],
		m.leds[nLEDs-1], m.leds[4],
		m.leds[nLEDs-2], m.leds[5],
		m.leds[nLEDs-3], m.leds[6],
		m.leds[nLEDs-4], m.leds[7],
		m.leds[nLEDs-5], m.leds[8],
		m.leds[nLEDs-6], m.leds[9],
		m.leds[nLEDs-7], m.leds[10],
		m.leds[nLEDs-8], m.leds[11],
		m.leds[15], m.leds[14], m.leds[13], m.leds[12],
	)
}
