package service

import (
	"fmt"
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/repository"
	"github.com/francescarpi/mytime/internal/types"
	"github.com/francescarpi/mytime/internal/util"
)

type WorkedDurationFormatted struct {
	Daily          string
	DailyGoal      string
	DailyOvertime  string
	Weekly         string
	WeeklyGoal     string
	WeeklyOvertime string
}

type Service struct {
	Repo repository.Repository
}

type SummaryDuration struct {
	Reported    string
	NotReported string
}

func (s *Service) GetTasksByDate(date time.Time) ([]model.Task, error) {
	tasks, err := s.Repo.GetTasksByDate(date)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) GetWorkedDuration(date time.Time) (WorkedDurationFormatted, error) {
	var result WorkedDurationFormatted

	daily, err := s.Repo.GetWorkedDurationForDate(date, types.All)
	if err != nil {
		return result, err
	}

	weekly, err := s.Repo.GetWeeklyWorkedDurationForDate(date)
	if err != nil {
		return result, err
	}

	settings, err := s.Repo.GetSettings()
	if err != nil {
		return result, err
	}

	dailyGoalSeconds := settings.GoalDayInSeconds(date)
	weeklyGoalSeconds := settings.GoalWeekInSeconds()

	result.Daily = util.HumanizeDuration(daily)
	result.DailyGoal = util.HumanizeDuration(dailyGoalSeconds)
	result.Weekly = util.HumanizeDuration(weekly)
	result.WeeklyGoal = util.HumanizeDuration(weeklyGoalSeconds)
	result.DailyOvertime = util.HumanizeDuration(daily - dailyGoalSeconds)
	result.WeeklyOvertime = util.HumanizeDuration(weekly - weeklyGoalSeconds)

	return result, nil
}

func (s *Service) GetSummaryDuration(date time.Time) (SummaryDuration, error) {
	var result SummaryDuration

	reported, err := s.Repo.GetWorkedDurationForDate(date, types.Reported)
	if err != nil {
		return result, err
	}

	notReported, err := s.Repo.GetWorkedDurationForDate(date, types.NotReported)
	if err != nil {
		return result, err
	}

	rawReported := float64(reported) / 3600
	rawNotReported := float64(notReported) / 3600

	result.Reported = fmt.Sprintf("%s (%.2f)", util.HumanizeDuration(reported), rawReported)
	result.NotReported = fmt.Sprintf("%s (%.2f)", util.HumanizeDuration(notReported), rawNotReported)

	return result, nil
}

func (s *Service) CreateTask(description string, project, externalId *string) error {
	s.Repo.CloseOpenedTasks()
	if err := s.Repo.CreateTask(description, project, externalId); err != nil {
		return err
	}
	return nil
}

func (s *Service) StartStopTask(id uint) error {
	task, err := s.Repo.GetTask(id)
	if err != nil {
		return err
	}

	if task.IsOpen() {
		return s.Repo.CloseTask(id)
	}

	return s.CreateTask(task.Desc, task.Project, task.ExternalId)
}

func (s *Service) UpdateTask(task *model.Task) error {
	if err := s.Repo.UpdateTask(task); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteTask(id uint) error {
	if err := s.Repo.DeleteTask(id); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetTasksToSync() []types.TasksToSync {
	tasks, err := s.Repo.GetTasksToSync()
	if err != nil {
		return []types.TasksToSync{}
	}
	return tasks
}

func (s *Service) SetTaskAsReported(id uint) error {
	return s.Repo.SetTaskAsReported(id)
}
