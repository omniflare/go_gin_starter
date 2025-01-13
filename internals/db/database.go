package db

import (
	"github.com/omniflare/go_starter/internals/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func Init(config *config.Config, DBMigrator func(*gorm.DB) error) *gorm.DB {
	uri := config.Database
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err.Error())
	}

	log.Printf("Connected to the database")

	if err := DBMigrator(db); err != nil {
		log.Fatalf("Unable to migrate the database: %v", err)
	}

	return db
}
