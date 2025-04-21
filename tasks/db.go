package tasks

import (
	"fmt"
	"log"
	"mytime/utils"
	"runtime/debug"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const DURATION = "COALESCE(STRFTIME('%s', end), STRFTIME('%s', DATETIME('now', 'localtime'))) - STRFTIME('%s', start)"
const ORDER = "start DESC, id"

type TasksManager struct {
	Conn *gorm.DB
}

func (t *TasksManager) GetTasksByDate(date string) ([]Task, error) {
	var tasks []Task
	t.Conn.
		Select(fmt.Sprintf("*, %s AS duration", DURATION)).
		Where("DATE(start) = DATE(?)", date).
		Order(ORDER).
		Find(&tasks)

	return tasks, nil
}

func (t *TasksManager) GetWorkedDaily(date string) (float64, error) {
	var result Duration
	t.Conn.
		Model(&Task{}).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0) AS duration", DURATION)).
		Where("DATE(start) = DATE(?)", date).
		Find(&result)

	return result.Duration, nil
}

func (t *TasksManager) GetWorkedWeekly(date string) (float64, error) {
	var result Duration
	dateParsed, err := utils.ParseDate(date)

	if err != nil {
		return 0, err
	}

	t.Conn.
		Model(&Task{}).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0) AS duration", DURATION)).
		Where("STRFTIME('%V', start) = ? AND STRFTIME('%Y', start) = ?", strconv.Itoa(dateParsed.Week), strconv.Itoa(dateParsed.Year)).
		Find(&result)

	return result.Duration, nil
}

func (t *TasksManager) GetTasksToSync() []TasksToSync {
	var result []TasksToSync

	t.Conn.
		Model(&Task{}).
		Select(fmt.Sprintf("GROUP_CONCAT(id, '-') as id, "+
			"external_id, "+
			"SUM(%s) AS duration, "+
			"desc, "+
			"STRFTIME('%%Y-%%m-%%d', start) AS date, "+
			"project, "+
			"GROUP_CONCAT(id) AS ids", DURATION)).
		Where("end IS NOT NULL AND reported = false AND external_id IS NOT NULL AND external_id != ''").
		Group("external_id, desc, STRFTIME('%Y-%m-%d', start), project").
		Order("STRFTIME('%Y-%m-%d', start) DESC, id").
		Find(&result)

	return result
}

func (t *TasksManager) StartStopTask(ID uint) {
	var task Task
	t.Conn.First(&task, ID)

	if task.ID == 0 {
		return
	}

	now := time.Now()
	if task.End == nil {
		task.End = &NaiveTime{now}
		t.Conn.Save(&task)
	} else {
		t.CloseOpenedTasks()
		newTask := Task{
			Desc:       task.Desc,
			Start:      NaiveTime{now},
			ExternalId: task.ExternalId,
			Project:    task.Project,
		}
		if err := t.Conn.Save(&newTask).Error; err != nil {
			log.Printf("Error saving new task: %v\nStacktrace:\n%s", err, debug.Stack())
		}
	}
}

func (t *TasksManager) DuplicateTaskWithDescription(ID uint, description string) {
	var task Task
	t.Conn.First(&task, ID)

	if task.ID == 0 {
		return
	}

	t.CloseOpenedTasks()

	// If the task is opened, we need to close it
	if task.End == nil {
		task.End = &NaiveTime{time.Now()}
		t.Conn.Save(&task)
	}

	now := time.Now()
	newTask := Task{
		Desc:       description,
		Start:      NaiveTime{now},
		ExternalId: task.ExternalId,
		Project:    task.Project,
	}

	if err := t.Conn.Save(&newTask).Error; err != nil {
		log.Printf("Error saving new task: %v\nStacktrace:\n%s", err, debug.Stack())
	}
}

func (t *TasksManager) CloseOpenedTasks() {
	var tasks []Task
	t.Conn.Where("end IS NULL").Find(&tasks)

	for _, task := range tasks {
		task.End = &NaiveTime{time.Now()}
		if err := t.Conn.Save(&task).Error; err != nil {
			log.Printf("Error saving new task: %v\nStacktrace:\n%s", err, debug.Stack())
		}
	}
}

func (t *TasksManager) GetTaskById(ID uint) (Task, error) {
	var task Task
	if err := t.Conn.First(&task, ID).Error; err != nil {
		return Task{}, err
	}
	return task, nil
}

func (t *TasksManager) Update(ID uint, project, description, externalId, start, end string) {
	var task Task
	t.Conn.First(&task, ID)

	if task.ID == 0 {
		return
	}

	task.Desc = description
	task.Start = NaiveTime{utils.UpdateTime(&task.Start.Time, start)}

	if project != "" {
		task.Project = &project
	}

	if externalId != "" {
		task.ExternalId = &externalId
	}

	if end != "" {
		task.End = &NaiveTime{utils.UpdateTime(&task.Start.Time, end)}
	}

	t.Conn.Save(&task)
}

func (t *TasksManager) Delete(ID uint) {
	var task Task
	t.Conn.First(&task, ID)

	if task.ID == 0 {
		return
	}

	if err := t.Conn.Delete(&task).Error; err != nil {
		log.Printf("Error deleting task: %v\nStacktrace:\n%s", err, debug.Stack())
	}
}

func (t *TasksManager) New(project, description, extID string) {
	now := time.Now()
	newTask := Task{
		Project:    &project,
		Desc:       description,
		ExternalId: &extID,
		Start:      NaiveTime{now},
	}

	if err := t.Conn.Save(&newTask).Error; err != nil {
		log.Printf("Error saving new task: %v\nStacktrace:\n%s", err, debug.Stack())
	}
}
