package pdb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex"`
	PublicKey []byte
}

func NewUser(username string, publicKey []byte) *User {
	return &User{
		Username:  username,
		PublicKey: publicKey,
	}
}
