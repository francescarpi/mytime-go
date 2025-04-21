package utils

import (
	"fmt"
	"math"
	"time"
)

type DateParsed struct {
	Year  int
	Month int
	Day   int
	Week  int
}

func ParseDate(date string) (DateParsed, error) {
	dateParsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing the date", err)
		return DateParsed{}, err
	}
	_, week := dateParsed.ISOWeek()
	return DateParsed{
		Year:  dateParsed.Year(),
		Month: int(dateParsed.Month()),
		Day:   dateParsed.Day(),
		Week:  week,
	}, nil
}

func HumanizeDuration(seconds float64) string {
	hours := math.Floor(seconds / 3600.0)
	minutes := math.Floor(math.Mod(seconds, 3600.0) / 60.0)
	return fmt.Sprintf("%dh%dm", int(hours), int(minutes))
}

func UpdateTime(date *time.Time, hour string) time.Time {
	parsed, err := time.Parse("15:04", hour)
	if err != nil {
		fmt.Println("Error parsing the hour", err)
		return *date
	}
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour(), parsed.Minute(), 0, 0, date.Location())
}
