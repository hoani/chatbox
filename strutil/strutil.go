package strutil

import "strings"

func SplitWidth(input string, width int) []string {
	words := strings.Fields(input)
	lines := []string{}
	var line string
	for _, word := range words {
		if len(line)+len(word)+1 > width {
			if line != "" {
				lines = append(lines, line)
			}
			for len(word) > width {
				lines = append(lines, word[:width])
				word = word[width:]
			}
			line = word
			continue
		}
		if line != "" {
			line += " "
		}
		line += word
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
