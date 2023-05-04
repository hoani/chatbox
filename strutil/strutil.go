package strutil

import (
	"strings"
)

func Simplify(input string) string {

	var vowels = map[rune]struct{}{
		'a': {}, 'e': {}, 'i': {}, 'o': {}, 'u': {},
		'A': {}, 'E': {}, 'I': {}, 'O': {}, 'U': {},
	}

	runes := []rune(input)

	trimRepeatedPrefix := func(runes []rune) []rune {
		start := runes[0]
		var count int
		var r rune
		for _, r = range runes {
			if r != start {
				break
			}
			count++
		}
		if count < 3 {
			return runes
		}
		if _, ok := vowels[start]; ok {
			return runes[count-2:]
		}
		return runes[count-1:]
	}

	index := 0
	for index < len(runes) {
		runes = append(runes[:index], trimRepeatedPrefix(runes[index:])...)
		index++
	}

	return string(runes)
}

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
		case "(", "[", "{", "<":
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
		case ")", "]", "}", ">":
			current += string(c)
			if depth > 0 && braceMap[token] == brace {
				depth--
				if depth == 0 {
					result = append(result, strings.TrimSpace(current))
					current = ""
				}
			}
		case "*":
			if depth == 0 {
				brace = token
				if current := strings.TrimSpace(current); current != "" {
					result = append(result, current)
				}
				current = token
				depth = 1
			} else if depth == 1 && brace == token {
				current += token
				result = append(result, strings.TrimSpace(current))
				current = ""
				depth = 0
			} else {
				current += token
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
