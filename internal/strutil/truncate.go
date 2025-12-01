package strutil

import "fmt"

// TruncateMiddle metni baştan ve sondan 'n' karakter alarak kısaltır.
func TruncateMiddle(text string, n int) string {
	// 1. Önce Rune'a çeviriyoruz (Emoji ve Türkçe karakter güvenliği için)
	runes := []rune(text)
	l := len(runes)

	// 2. Güvenlik Kontrolü: Eğer metin zaten kısaysa olduğu gibi döndür
	// "n * 2" çünkü baştan n, sondan n alacağız.
	if l <= n*2 {
		return text
	}

	// 3. Birleştirme: Baş + ... + Son
	// string(runes[...]) dönüşümü ile byte hatası yapmadan birleştiriyoruz
	return fmt.Sprintf("%s...%s", string(runes[:n]), string(runes[l-n:]))
}
