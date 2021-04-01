package pdb

type NotificationType int
type NotificationState int

type Notification struct {
	SenderID   uint `gorm:"primaryKey;autoincrement:false"`
	ReceiverID uint `gorm:"primaryKey;autoincrement:false"`
	Type       NotificationType
	State      NotificationState
	Sender     User `gorm:"foreignKey:SenderID"`
	Receiver   User `gorm:"foreignKey:ReceiverID"`
}

func NewNotification(ntype NotificationType, sender *User, receiver *User, state NotificationState) *Notification {
	return &Notification{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Type:       ntype,
		State:      state,
	}
}
