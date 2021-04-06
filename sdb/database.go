package sdb

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateConnection(dsn string, log *log.Logger) (*gorm.DB, error) {
	newLogger := logger.New(
		log,
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      false,
		},
	)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Request{})
	db.AutoMigrate(&Friend{})
	db.AutoMigrate(&Message{})

	return db, nil
}
