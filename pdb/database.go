package pdb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateConnection(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err

	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Notification{})
	db.AutoMigrate(&Message{})
	return db, nil
}
