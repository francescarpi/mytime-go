package settings

import (
	"strconv"
	"strings"
	"time"
)

type Settings struct {
	ID                uint    `gorm:"primarykey"`
	Integration       *string `gorm:"type:varchar"`
	WorkHours         string  `gorm:"not null;type:varchar"`
	Theme             string  `gorm:"not null;type:varchar"`
	ViewType          string  `gorm:"not null;type:varchar"`
	DarkMode          bool    `gorm:"not null;type:boolean"`
	RightSideBarOpen  bool    `gorm:"not null;column:right_sidebar_open;default:false"`
	ThemeSecondary    string  `gorm:"not null;type:VARCHAR;default:#ce93d8"`
	IntegrationConfig string  `gorm:"not null;type:text;default:{}"`
}

func (s *Settings) GoalDayInSeconds(date time.Time) float64 {
	// parts are the work hours for each day of the week, where the first element is Monday
	parts := strings.Split(s.WorkHours, ",")

	// Weekday returns 0 for Sunday, 1 for Monday, ..., 6 for Saturday
	weekday := int(date.Weekday()) - 1
	if weekday < 0 {
		weekday = 6
	}
	goalDayStr := parts[weekday]
	goalDay, err := strconv.ParseFloat(goalDayStr, 64)

	if err != nil {
		return 0
	}
	return goalDay * 3600.0
}

func (s *Settings) GoalWeekInSeconds() float64 {
	parts := strings.Split(s.WorkHours, ",")
	goalWeek := 0.0
	for _, part := range parts {
		goalDay, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return 0
		}
		goalWeek += goalDay
	}
	return goalWeek * 3600.0
}
