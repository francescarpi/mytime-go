package service

import (
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/repository"
	"github.com/francescarpi/mytime/internal/util"
)

type WorkedDuration struct {
	DailyFormatted          string
	DailyGoalFormatted      string
	DailyOvertime           string
	WeeklyFormatted         string
	WeeklyGoalFormatted     string
	WeeklyOvertimeFormatted string
}

type Service struct {
	Repo repository.Repository
}

func (s *Service) GetTasksByDate(date time.Time) ([]model.Task, error) {
	tasks, err := s.Repo.GetTasksByDate(date)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) GetWorkedDuration(date time.Time) (WorkedDuration, error) {
	var result WorkedDuration

	daily, err := s.Repo.GetWorkedDurationForDate(date)
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

	result.DailyFormatted = util.HumanizeDuration(daily)
	result.DailyGoalFormatted = util.HumanizeDuration(dailyGoalSeconds)
	result.WeeklyFormatted = util.HumanizeDuration(weekly)
	result.WeeklyGoalFormatted = util.HumanizeDuration(weeklyGoalSeconds)

	if daily > dailyGoalSeconds {
		result.DailyOvertime = util.HumanizeDuration(daily - dailyGoalSeconds)
	}

	if weekly > weeklyGoalSeconds {
		result.WeeklyOvertimeFormatted = util.HumanizeDuration(weekly - weeklyGoalSeconds)
	}

	return result, nil
}

func (s *Service) CreateTask(description string, project, externalId *string) error {
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

	return s.Repo.CreateTask(task.Desc, task.Project, task.ExternalId)
}

func (s *Service) UpdateTask(task *model.Task) error {
	if err := s.Repo.UpdateTask(task); err != nil {
		return err
	}
	return nil
}
