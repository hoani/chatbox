package strutil

import "strings"

func SplitBrackets(input string) []string {
	result := []string{}
	current := ""
	depth := 0
	brace := ""
	braceMap := map[string]string{
		")": "(",
		"]": "[",
		"}": "{",
		">": "<",
		"*": "*",
	}

	for _, c := range input {
		token := string(c)
		switch token {
		case "(", "[", "{", "<", "*":
			if depth == 0 {
				brace = token
				if current := strings.TrimSpace(current); current != "" {
					result = append(result, current)
				}
				current = ""
			}
			current += token
			if token == brace {
				depth++
			}
		case ")", "]", "}", ">", "*":
			current += string(c)
			if depth > 0 && braceMap[token] == brace {
				depth--
				if depth == 0 {
					result = append(result, strings.TrimSpace(current))
					current = ""
				}
			}
		default:
			current += string(c)
		}
	}
	if current != "" {
		return append(result, strings.TrimSpace(current))
	}
	return result
}

func SplitSentences(input string) []string {
	result := []string{}
	current := ""
	punctuation := ""
	for _, c := range input {
		token := string(c)
		switch token {
		case ".", "!", "?":
			if punctuation == "" {
				punctuation = token
			}
			current += token
		default:
			if punctuation != "" {
				result = append(result, strings.TrimSpace(current))
				current = ""
				punctuation = ""
			}
			current += token
		}
	}
	if current != "" {
		return append(result, strings.TrimSpace(current))
	}
	return result
}

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
