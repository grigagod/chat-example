package pdb

type User struct {
	Username  string `gorm:"primaryKey"`
	PublicKey []byte
}

func NewUser(username string, publicKey []byte) *User {
	return &User{
		Username:  username,
		PublicKey: publicKey,
	}
}
