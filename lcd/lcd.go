package lcd

import (
	"os/exec"
	"strings"
)

const (
	Width  = 16
	Height = 2
)

type cmd []string

func Cmd() cmd {
	return make(cmd, 0)
}

func (c cmd) Line1(in string) cmd {
	return append(c, "--line1", Whitespace(in))
}

func (c cmd) Line2(in string) cmd {
	return append(c, "--line2", Whitespace(in))
}

func (c cmd) RGB(in string) cmd {
	return append(c, "--rgb", in)
}

func (c cmd) Init() cmd {
	return append(c, "--init", "true")
}

func (c cmd) Run() {
	args := []string{"lcd/lcd.py"}
	for _, v := range c {
		args = append(args, v)
	}
	exec.Command("python3", args...).Run()
}

func Pad(in string) string {
	if len(in) > Width {
		return in[:Width]
	}
	padding := (Width - len(in)) / 2
	return strings.Repeat(" ", padding) + in + strings.Repeat(" ", padding)
}

func Whitespace(in string) string {
	if diff := Width - len(in); diff > 0 {
		return in + strings.Repeat(" ", diff)
	}
	return in
}
