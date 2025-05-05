package util

import (
	"fmt"
)

func HumanizeDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if minutes < 0 {
		minutes = -minutes
	}

	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, minutes)
}
