package util

import "testing"

func TestHumanizeDuration(t *testing.T) {
	tests := []struct {
		duration int
		expected string
	}{
		{0, "0m"},
		{1, "0m"},
		{2, "0m"},
		{60, "1m"},
		{61, "1m"},
		{120, "2m"},
		{3600, "1h"},
		{3661, "1h1m"},
		{-3661, "-1h1m"},
	}

	for _, test := range tests {
		result := HumanizeDuration(test.duration)
		if result != test.expected {
			t.Errorf("Expected %v for duration %v, got %v", test.expected, test.duration, result)
		}
	}
}
