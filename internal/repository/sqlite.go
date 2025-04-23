package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DURATION = "COALESCE(STRFTIME('%s', end), STRFTIME('%s', DATETIME('now', 'localtime'))) - STRFTIME('%s', start)"
const ORDER = "start DESC, id"

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepository(dsn string) *SqliteRepository {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Task{})
	db.AutoMigrate(&model.Settings{})

	return &SqliteRepository{db: db}
}

func (r *SqliteRepository) GetTasksByDate(date time.Time) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.
		Select(fmt.Sprintf("*, %s AS duration", DURATION)).
		Where("DATE(start) = DATE(?)", date.Format(time.DateOnly)).
		Order(ORDER).
		Find(&tasks).
		Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *SqliteRepository) GetWorkedDurationForDate(date time.Time) (int, error) {
	var result int
	err := r.db.
		Model(&model.Task{}).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0) AS duration", DURATION)).
		Where("DATE(start) = DATE(?)", date.Format(time.DateOnly)).
		Find(&result).
		Error

	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *SqliteRepository) GetWeeklyWorkedDurationForDate(date time.Time) (int, error) {
	var result int
	_, weekNumber := date.ISOWeek()
	err := r.db.
		Model(&model.Task{}).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0) AS duration", DURATION)).
		Where("STRFTIME('%V', start) = ? AND STRFTIME('%Y', start) = ?", strconv.Itoa(weekNumber), strconv.Itoa(date.Year())).
		Find(&result).
		Error

	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *SqliteRepository) GetSettings() (*model.Settings, error) {
	var settings model.Settings
	if err := r.db.First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *SqliteRepository) CreateTask(description string, project, externalId *string) error {
	r.CloseOpenedTasks()

	newTask := model.Task{
		Project:    project,
		Desc:       description,
		ExternalId: externalId,
		Start:      model.LocalTimestamp{Time: time.Now()},
	}

	if err := r.db.Save(&newTask).Error; err != nil {
		return err
	}
	return nil
}

func (r *SqliteRepository) CloseOpenedTasks() error {
	var tasks []model.Task
	r.db.Where("end IS NULL").Find(&tasks)

	for _, task := range tasks {
		task.End = &model.LocalTimestamp{Time: time.Now()}
		if err := r.db.Save(&task).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *SqliteRepository) CloseTask(id uint) error {
	var task model.Task
	err := r.db.First(&task, id).Error

	if err != nil {
		return err
	}

	if task.End != nil {
		return fmt.Errorf("task already closed")
	}

	task.End = &model.LocalTimestamp{Time: time.Now()}
	if err := r.db.Save(&task).Error; err != nil {
		return err
	}

	return nil
}

func (r *SqliteRepository) GetTask(id uint) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, id).Error

	if err != nil {
		return nil, err
	}

	return &task, nil
}
