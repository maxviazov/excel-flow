package textutil

import (
	"html"
	"strings"
	"unicode"
)

// CleanText удаляет HTML entities и очищает текст
func CleanText(s string) string {
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	return strings.TrimSpace(s)
}

// TransliterateToHebrew транслитерирует латинские буквы в иврит
func TransliterateToHebrew(s string) string {
	if containsHebrew(s) {
		return s
	}
	
	translitMap := map[rune]string{
		'a': "א", 'A': "א",
		'b': "ב", 'B': "ב",
		'c': "ק", 'C': "ק",
		'd': "ד", 'D': "ד",
		'e': "א", 'E': "א",
		'f': "פ", 'F': "פ",
		'g': "ג", 'G': "ג",
		'h': "ה", 'H': "ה",
		'i': "י", 'I': "י",
		'j': "ג", 'J': "ג",
		'k': "ק", 'K': "ק",
		'l': "ל", 'L': "ל",
		'm': "מ", 'M': "מ",
		'n': "נ", 'N': "נ",
		'o': "ו", 'O': "ו",
		'p': "פ", 'P': "פ",
		'q': "ק", 'Q': "ק",
		'r': "ר", 'R': "ר",
		's': "ס", 'S': "ס",
		't': "ט", 'T': "ט",
		'u': "ו", 'U': "ו",
		'v': "ו", 'V': "ו",
		'w': "ו", 'W': "ו",
		'x': "קס", 'X': "קס",
		'y': "י", 'Y': "י",
		'z': "ז", 'Z': "ז",
	}
	
	var result strings.Builder
	for _, r := range s {
		if heb, ok := translitMap[r]; ok {
			result.WriteString(heb)
		} else {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func containsHebrew(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Hebrew, r) {
			return true
		}
	}
	return false
}

// SanitizeForMOH очищает текст для отправки в систему Министерства здравоохранения
func SanitizeForMOH(s string) string {
	s = CleanText(s)
	s = TransliterateToHebrew(s)
	return s
}
