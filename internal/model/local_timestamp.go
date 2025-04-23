package model

import (
	"database/sql/driver"
	"time"
)

type LocalTimestamp struct {
	time.Time
}

func (nt LocalTimestamp) Value() (driver.Value, error) {
	return nt.Format("2006-01-02 15:04:05.999999"), nil
}

func (nt *LocalTimestamp) Scan(value any) error {
	if value == nil {
		return nil
	}
	nt.Time = value.(time.Time)
	return nil
}
