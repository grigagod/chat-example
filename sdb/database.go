package sdb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Request{})
	db.AutoMigrate(&Friend{})
	db.AutoMigrate(&Message{})

	return db, nil
}






