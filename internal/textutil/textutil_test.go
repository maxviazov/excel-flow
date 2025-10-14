package textutil

import "testing"

func TestCleanText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"סופר פאפא - סניף רח&#39; שבי ציון", "סופר פאפא - סניף רח' שבי ציון"},
		{"בע&quot;מ - סניף", "בע\"מ - סניף"},
		{"רגיל טקסט", "רגיל טקסט"},
	}

	for _, tt := range tests {
		result := CleanText(tt.input)
		if result != tt.expected {
			t.Errorf("CleanText(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestTransliterateToHebrew(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Maximil", "מאקסימיל"},
		{"מקסמיל", "מקסמיל"}, // уже иврит, не трогаем
		{"Test", "טאסט"},
	}

	for _, tt := range tests {
		result := TransliterateToHebrew(tt.input)
		if result != tt.expected {
			t.Errorf("TransliterateToHebrew(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeForMOH(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"סופר פאפא - סניף רח&#39; שבי ציון", "סופר פאפא - סניף רח' שבי ציון"},
		{"Maximil", "מאקסימיל"},
		{"בע&quot;מ", "בע\"מ"},
	}

	for _, tt := range tests {
		result := SanitizeForMOH(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeForMOH(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
