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
	SenderName   string `gorm:"primaryKey;autoincrement:false"`
	ReceiverName string `gorm:"primaryKey;autoincrement:false"`
	State        NotificationState
	Sender       User `gorm:"foreignKey:SenderName"`
	Receiver     User `gorm:"foreignKey:ReceiverName"`
}

func NewNotification(sender *User, receiver *User) *Notification {
	return &Notification{
		SenderName:   sender.Username,
		ReceiverName: receiver.Username,
		State:        Initiated,
	}
}
