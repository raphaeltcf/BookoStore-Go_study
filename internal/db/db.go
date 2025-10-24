package db

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	once sync.Once
)

func Connect(dsn string, models ...interface{}) *gorm.DB {
	once.Do(func() {
		var err error
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		if err := DB.AutoMigrate(models...); err != nil {
			log.Fatalf("Failed to migrate models: %v", err)
		}
	})
	return DB
}