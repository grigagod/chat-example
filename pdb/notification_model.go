package pdb

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
	SenderName   string `gorm:"primaryKey"`
	ReceiverName string `gorm:"primaryKey"`
	State        NotificationState
	Sender       User `gorm:"foreignKey:SenderName;References:Username"`
	Receiver     User `gorm:"foreignKey:ReceiverName;References:Username"`
}

func NewNotification(sender *User, receiver *User) *Notification {
	return &Notification{
		SenderName:   sender.Username,
		ReceiverName: receiver.Username,
		State:        Initiated,
	}
}
