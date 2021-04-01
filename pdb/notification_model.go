package pdb

type NotificationType int

const (
	KeyExchange NotificationType = iota
)

type NotificationState int

const (
	//Initiated means some notification initiated by sender
	Initiated NotificationState = iota
	//Received means some notification received by receiver
	Received
	//Accepted means some notification accepted by receiver
	Accepted
	//Declined means some notification declined by receiver
	Declined
)

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
