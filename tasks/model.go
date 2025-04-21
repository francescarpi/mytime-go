package tasks

import (
	"database/sql/driver"
	"log"
	"mytime/utils"
	"time"
)

type NaiveTime struct {
	time.Time
}

func (nt NaiveTime) Value() (driver.Value, error) {
	return nt.Format("2006-01-02 15:04:05.999999"), nil
}

func (nt *NaiveTime) Scan(value any) error {
	log.Println("Scanning value:", value)
	if value == nil {
		return nil
	}
	nt.Time = value.(time.Time)
	return nil
}

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
