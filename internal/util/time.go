package util

import "time"

func UpdateTime(date *time.Time, hour string) (time.Time, error) {
	parsed, err := time.Parse("15:04", hour)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		date.Location(),
	), nil
}
