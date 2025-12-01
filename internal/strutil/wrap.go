package strutil

import "strings"

// HardWrap wraps the text at the specified limit.
func HardWrap(text string, limit int) string {
	if len(text) <= limit {
		return text
	}

	var chunks []string
	// Using rune slice prevents corruption in Unicode characters (emojis etc.).
	// Base64 is usually ASCII but using runes is a good habit.
	runes := []rune(text)

	for i := 0; i < len(runes); i += limit {
		end := i + limit
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return strings.Join(chunks, "\n")
}
