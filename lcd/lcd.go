package lcd

import "strings"

const (
	Width  = 16
	Height = 2
)

func Pad(in string) string {
	if len(in) > Width {
		return in[:Width]
	}
	padding := (Width - len(in)) / 2
	return strings.Repeat(" ", padding) + in + strings.Repeat(" ", padding)
}
