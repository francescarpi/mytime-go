package repository

import (
	"time"

	"github.com/francescarpi/mytime/internal/model"
)

type Repository interface {
	GetTasksByDate(date time.Time) ([]model.Task, error)
	GetWorkedDurationForDate(date time.Time) (int, error)
	GetWeeklyWorkedDurationForDate(date time.Time) (int, error)
	GetSettings() (*model.Settings, error)
}
