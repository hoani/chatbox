package hal

import (
	"os/exec"
)

type lcd struct{
	line1 string
	line2 string
	color LCDColor
}

func newLCD() *lcd {
	return &lcd{
		color: LCDRed,
	}
}

func (l *lcd) Write(line1, line2 string, color LCDColor) {
	args := []string{"lcd/lcd.py"}
	if l.line1 != line1 {
		args = append(args, "--line1", line1)
		l.line1 = line1
	}
	if l.line2 != line2 {
		args = append(args, "--line2", line2)
		l.line2 = line2
	}
	if l.color != color {
		l.color = color
	}
	if len(args) > 1 {
		args = append(args, "--rgb", LCDColorToString(l.color))
		exec.Command("python3", args...).Run()
	}
}

func LCDColorToString(color LCDColor) string {
	switch color {
	case LCDRed:
		return "255,0,0"
	case LCDGreen:
		return "255,209,0"
	case LCDAqua:
		return "248,248,60"
	case LCDBlue:
		return "255,255,255"
	default:
		return "255,255,255"
	}
}
