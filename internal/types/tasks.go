package types

import (
	"database/sql/driver"
	"strings"
)

type ListOfIds struct {
	IDs []string
}

func (i ListOfIds) Value() (driver.Value, error) {
	return "", nil
}

func (i *ListOfIds) Scan(value any) error {
	if value == nil {
		return nil
	}

	i.IDs = strings.Split(value.(string), ",")
	return nil
}

type TasksToSync struct {
	Id         string
	ExternalId string
	Duration   int
	Desc       string
	Date       string
	Project    string
	Ids        ListOfIds
}

type TaskStatus int64

const (
	Reported TaskStatus = iota
	NotReported
	All
)
