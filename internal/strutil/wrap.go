package strutil

import "strings"

func HardWrap(text string, limit int) string {
	if len(text) <= limit {
		return text
	}

	var chunks []string
	// Rune slice kullanarak Unicode karakterlerde (emoji vb.) bozulmayı önleriz.
	// Base64 genelde ASCII'dir ama alışkanlık olarak rune kullanmak iyidir.
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
