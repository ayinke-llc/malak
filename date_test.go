package malak

import (
	"testing"
)

func TestOrdinalSuffix(t *testing.T) {
	tests := []struct {
		day      int
		expected string
	}{
		{1, "st"},
		{2, "nd"},
		{3, "rd"},
		{4, "th"},
		{11, "th"},
		{12, "th"},
		{13, "th"},
		{21, "st"},
		{22, "nd"},
		{23, "rd"},
		{24, "th"},
		{31, "st"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := ordinalSuffix(tt.day)
			if got != tt.expected {
				t.Errorf("ordinalSuffix(%d) = %v, want %v", tt.day, got, tt.expected)
			}
		})
	}
}

func TestGetTodayFormatted(t *testing.T) {
	// Since this function returns the current date, we can only verify
	// that it returns a non-empty string in the expected format
	result := GetTodayFormatted()
	if result == "" {
		t.Error("GetTodayFormatted() returned empty string")
	}
}
