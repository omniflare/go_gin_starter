package db

import (
	models "github.com/omniflare/go_starter/internals/models/user"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
