package db

import (
	"mytime/config"
	"mytime/settings"
	"mytime/tasks"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetConnection(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&tasks.Task{})
	db.AutoMigrate(&settings.Settings{})

	return db
}
