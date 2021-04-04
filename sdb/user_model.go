package sdb

type User struct {
	Username   string `gorm:"primaryKey"`
	PrivateKey []byte `gorm:"unique"`
}

func NewUser(username string, privateKey []byte) *User {
	return &User{
		Username:   username,
		PrivateKey: privateKey,
	}
}

