package tasks

import (
	"mytime/utils"
	"time"
)

type Task struct {
	ID         uint       `gorm:"primarykey"`
	Desc       string     `gorm:"not null;type:varchar"`
	Start      NaiveTime  `gorm:"not null;type:timestamp"`
	End        *NaiveTime `gorm:"default:NULL;type:timestamp"`
	Reported   bool       `gorm:"default:false"`
	ExternalId *string    `gorm:"default:NULL;type:varchar"`
	Project    *string    `gorm:"default:NULL;type:varchar"`
	Favourite  bool       `gorm:"default:false"`
	Duration   float64    `gorm:"-:migration;->"`
}

func (t *Task) GetProjectOrBlank() string {
	if t.Project == nil {
		return ""
	}
	return *t.Project
}

func (t *Task) GetExternalIdOrBlank() string {
	if t.ExternalId == nil {
		return ""
	}
	return *t.ExternalId
}

func (t *Task) GetEndFormatedOr(txt string) string {
	if t.End == nil {
		return txt
	}
	return t.End.Format(time.Kitchen)
}

func (t *Task) GetDurationHumanized() string {
	return utils.HumanizeDuration(t.Duration)
}

func (t *Task) GetReportedIcon() string {
	if t.Reported {
		return "✓"
	}
	return "✗"
}
