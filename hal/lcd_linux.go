package hal

import (
	"fmt"
	"os/exec"
)

type lcd struct{}

func newLCD() *lcd {
	return &lcd{}
}

func (l *lcd) Write(line1, line2 string, color LCDColor) {
	args := []string{"lcd/lcd.py", "--line1", line1, "--line2", line2, "--rgb", LCDColorToString(color)}
	exec.Command("python3", args...).Run()
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
		return "200,200,200"
	default:
		return "200,200,200"
	}
}
