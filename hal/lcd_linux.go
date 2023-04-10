package hal

import (
	"fmt"
	"os/exec"
)

type lcd struct {}

func newLCD() *lcd {
	return &lcd{}
}

func (l *lcd) Write(line1, line2 string, color *RGB) {
	args  := []string{"lcd/lcd.py", "--line1", line1, "--line2", line2}
	if color != nil {
		args = append(args, "--rgb", fmt.Sprintf("%d,%d,%d", color.R, color.G, color.B))
	}
	exec.Command("python3", args...).Run()
}
