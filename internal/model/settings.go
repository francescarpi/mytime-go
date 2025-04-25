package model

import (
	"fmt"
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

	// cache
	cachedHours      []float64
	cachedHoursValid bool
}

func (s *Settings) GoalDayInSeconds(date time.Time) int {
	hours, err := s.parsedWorkHours()
	if err != nil {
		return 0
	}
	weekday := int(date.Weekday()) - 1
	if weekday < 0 {
		weekday = 6
	}
	return int(hours[weekday] * 3600)
}

func (s *Settings) GoalWeekInSeconds() int {
	hours, err := s.parsedWorkHours()
	if err != nil {
		return 0
	}

	var total float64
	for _, h := range hours {
		total += h
	}

	return int(total * 3600)
}

func (s *Settings) parsedWorkHours() ([]float64, error) {
	if s.cachedHoursValid {
		return s.cachedHours, nil
	}

	parts := strings.Split(s.WorkHours, ",")
	if len(parts) != 7 {
		return nil, fmt.Errorf("invalid number of days in WorkHours")
	}
	hours := make([]float64, 7)
	for i, part := range parts {
		h, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, err
		}
		hours[i] = h
	}

	s.cachedHours = hours
	s.cachedHoursValid = true

	return hours, nil
}
