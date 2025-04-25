package util

import "fmt"

// Helper to apply conditional color
func Colorize(label, key string, enabled bool) string {
	keyColor := "blue"
	if !enabled {
		keyColor = "gray"
	}
	return fmt.Sprintf("[white]%s:[%s] %s [white]| ", label, keyColor, key)
}
