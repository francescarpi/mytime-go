package tasks

import (
	"database/sql/driver"
	"strings"
	"time"
)

type Duration struct {
	Duration float64
}

type IdList struct {
	Ids []string
}

func (i IdList) Value() (driver.Value, error) {
	return "", nil
}

func (i *IdList) Scan(value any) error {
	if value == nil {
		return nil
	}

	i.Ids = strings.Split(value.(string), ",")
	return nil
}

type TasksToSync struct {
	Id         string
	ExternalId string
	Duration   float64
	Desc       string
	Date       string
	Project    string
	Ids        IdList
}

type NaiveTime struct {
	time.Time
}

func (nt NaiveTime) Value() (driver.Value, error) {
	return nt.Format("2006-01-02 15:04:05.999999"), nil
}

func (nt *NaiveTime) Scan(value any) error {
	if value == nil {
		return nil
	}
	nt.Time = value.(time.Time)
	return nil
}
