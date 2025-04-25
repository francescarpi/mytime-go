package repository

import (
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/types"
)

type Repository interface {
	GetTasksByDate(date time.Time) ([]model.Task, error)
	GetTasksToSync() ([]types.TasksToSync, error)
	GetWorkedDurationForDate(date time.Time) (int, error)
	GetWeeklyWorkedDurationForDate(date time.Time) (int, error)
	GetSettings() (*model.Settings, error)
	CreateTask(description string, project, externalId *string) error
	CloseOpenedTasks() error
	CloseTask(id uint) error
	GetTask(id uint) (*model.Task, error)
	UpdateTask(task *model.Task) error
	DeleteTask(id uint) error
	SetTaskAsReported(id uint) error
}
