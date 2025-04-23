package model

type Task struct {
	ID         uint            `gorm:"primarykey"`
	Desc       string          `gorm:"not null;type:varchar"`
	Start      LocalTimestamp  `gorm:"not null;type:timestamp"`
	End        *LocalTimestamp `gorm:"default:NULL;type:timestamp"`
	Reported   bool            `gorm:"default:false"`
	ExternalId *string         `gorm:"default:NULL;type:varchar"`
	Project    *string         `gorm:"default:NULL;type:varchar"`
	Favourite  bool            `gorm:"default:false"`
	Duration   int             `gorm:"-:migration;->"`
}
