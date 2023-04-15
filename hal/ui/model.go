package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hoani/chatbox/strutil"
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

type Debug string

type LEDShow struct{}

type LEDClear struct{}

type LCDUpdate struct {
	Color *lipgloss.Color
	Lines [LCDHeight]string
}

type Model struct {
	pending     [nLEDs]lipgloss.Color
	leds        [nLEDs]string
	lcd         LCDUpdate
	lcdStyle    lipgloss.Style
	debug       string
	buttonState bool
	width       int
}

func NewUI() *Model {
	colors := [nLEDs]lipgloss.Color{}
	leds := [nLEDs]string{}
	for i := range colors {
		colors[i] = lipgloss.Color("#000000")
		leds[i] = lipgloss.NewStyle().Foreground(colors[i]).Render(" ")
	}
	lcdColor := lipgloss.Color("#0000ee")

	return &Model{
		pending: colors,
		leds:    leds,
		lcd: LCDUpdate{
			Color: &lcdColor,
			Lines: [LCDHeight]string{},
		},
		lcdStyle: lipgloss.NewStyle().
			ColorWhitespace(true).
			MaxWidth(4 + LCDWidth).PaddingLeft(2).PaddingRight(2).
			Background(lcdColor).
			Foreground(lipgloss.Color("#eeeeee")),
		width: 100,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateCount++
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeySpace:
			m.buttonState = !m.buttonState
			return m, nil
		}

	case LEDColors:
		for i, color := range msg.Colors {
			if i+msg.Index < nLEDs {
				m.pending[i] = color
			}
		}
	case LEDShow:
		for i, color := range m.pending {
			m.leds[i] = lipgloss.NewStyle().Foreground(color).Render("✮")
		}

	case LEDClear:
		for i := range m.pending {
			m.pending[i] = lipgloss.Color("#000000")
		}

	case LCDUpdate:
		m.lcd.Lines = msg.Lines
		for i, line := range m.lcd.Lines {
			if padding := LCDWidth - len(line); padding > 0 {
				m.lcd.Lines[i] = line + strings.Repeat(" ", padding)
			}
		}
		if msg.Color != nil {
			m.lcd.Color = msg.Color
			m.lcdStyle = m.lcdStyle.Background(m.lcd.Color)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case Debug:
		m.debug = string(msg)
	}
	return m, nil
}

func (m *Model) View() string {

	return m.ViewLEDs() + "\n" +
		helpStyle("Press space to talk/stop talking, press Ctl-C to quit UI")
}

func (m *Model) ViewLEDs() string {
	debugSplit := strutil.SplitWidth(m.debug, m.width)

	return fmt.Sprintf(`
      %s %s %s %s
    %s         %s
  %s             %s
%s                 %s
%s                 %s  %s
%s                 %s  %s
%s                 %s
  %s             %s
    %s         %s
      %s %s %s %s

%s
`, m.leds[0], m.leds[1], m.leds[2], m.leds[3],
		m.leds[nLEDs-1], m.leds[4],
		m.leds[nLEDs-2], m.leds[5],
		m.leds[nLEDs-3], m.leds[6],
		m.leds[nLEDs-4], m.leds[7], m.lcdStyle.Render(m.lcd.Lines[0]),
		m.leds[nLEDs-5], m.leds[8], m.lcdStyle.Render(m.lcd.Lines[1]),
		m.leds[nLEDs-6], m.leds[9],
		m.leds[nLEDs-7], m.leds[10],
		m.leds[nLEDs-8], m.leds[11],
		m.leds[15], m.leds[14], m.leds[13], m.leds[12],
		strings.Join(debugSplit, "\n"),
	)
}

func (m *Model) ButtonState() bool {
	return m.buttonState
}
