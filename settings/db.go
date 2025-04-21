package settings

import (
	"gorm.io/gorm"
)

func GetSettings(conn *gorm.DB) (*Settings, error) {
	var settings Settings
	if err := conn.First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}
