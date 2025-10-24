package textutil

import (
	"strings"
	"unicode"
)

// TransliterateToHebrew converts English text to Hebrew transliteration
func TransliterateToHebrew(text string) string {
	if text == "" {
		return text
	}
	
	// If already contains Hebrew, return as is
	for _, r := range text {
		if unicode.Is(unicode.Hebrew, r) {
			return text
		}
	}
	
	// English to Hebrew transliteration map
	translitMap := map[string]string{
		"a": "א", "A": "א",
		"b": "ב", "B": "ב",
		"c": "ק", "C": "ק",
		"d": "ד", "D": "ד",
		"e": "א", "E": "א",
		"f": "פ", "F": "פ",
		"g": "ג", "G": "ג",
		"h": "ה", "H": "ה",
		"i": "י", "I": "י",
		"j": "ג'", "J": "ג'",
		"k": "כ", "K": "כ",
		"l": "ל", "L": "ל",
		"m": "מ", "M": "מ",
		"n": "נ", "N": "נ",
		"o": "ו", "O": "ו",
		"p": "פ", "P": "פ",
		"q": "ק", "Q": "ק",
		"r": "ר", "R": "ר",
		"s": "ס", "S": "ס",
		"t": "ט", "T": "ט",
		"u": "ו", "U": "ו",
		"v": "ו", "V": "ו",
		"w": "ו", "W": "ו",
		"x": "קס", "X": "קס",
		"y": "י", "Y": "י",
		"z": "ז", "Z": "ז",
		"ch": "צ'", "Ch": "צ'", "CH": "צ'",
		"sh": "ש", "Sh": "ש", "SH": "ש",
		"th": "ת", "Th": "ת", "TH": "ת",
	}
	
	result := text
	
	// Replace multi-character combinations first
	for eng, heb := range translitMap {
		if len(eng) > 1 {
			result = strings.ReplaceAll(result, eng, heb)
		}
	}
	
	// Replace single characters
	var builder strings.Builder
	for _, r := range result {
		if heb, ok := translitMap[string(r)]; ok {
			builder.WriteString(heb)
		} else if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(r)
		}
	}
	
	return builder.String()
}
