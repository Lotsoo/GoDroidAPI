package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDb() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
