package strutil

import "fmt"

// TruncateMiddle truncates the text by taking 'n' characters from the start and end.
func TruncateMiddle(text string, n int) string {
	// 1. Convert to Rune slice (for Emoji and Unicode safety)
	runes := []rune(text)
	l := len(runes)

	// 2. Safety Check: If text is already short, return as is
	// "n * 2" because we take n from start and n from end.
	if l <= n*2 {
		return text
	}

	// 3. Concatenation: Start + ... + End
	// string(runes[...]) conversion prevents byte errors
	return fmt.Sprintf("%s...%s", string(runes[:n]), string(runes[l-n:]))
}
