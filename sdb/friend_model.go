package sdb

type Friend struct {
	FriendName string `gorm:"primaryKey"`
	SharedKey  []byte
	OwnerName  string `gorm:"primaryKey"`
	Owner      User   `gorm:"foreignKey:OwnerName;References:Username"`
}

func NewFriend(friendname string, sharedKey []byte, ownerName string) *Friend {
	return &Friend{
		FriendName: friendname,
		SharedKey:  sharedKey,
		OwnerName:  ownerName,
	}
}
