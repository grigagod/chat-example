package pdb

import "gorm.io/gorm"

type MessageState int

const (
	//Initiated means some notification initiated by sender
	MsgInitiated MessageState = iota
	//Received means some notification received by receiver
	MsgReceived
)

type Message struct {
	SenderName   string `gorm:"index"`
	ReceiverName string `gorm:"index"`
	State        MessageState
	Message      []byte
	Timestamp    int64
	Sender       User `gorm:"foreignKey:SenderName;References:Username"`
	Receiver     User `gorm:"foreignKey:ReceiverName;References:Username"`
}

func NewMessage(Sender *User, Receiver *User, timestamp int64, msg []byte) *Message {
	return &Message{
		SenderName:   Sender.Username,
		ReceiverName: Receiver.Username,
		State:        MsgInitiated,
		Timestamp:    timestamp,
		Message:      msg,
	}
}

func (msg *Message) UpdateMsgState(db *gorm.DB, newState MessageState) {
	db.Model(&Message{}).Where("sender_name = ? AND receiver_name = ? AND state = ? AND timestamp = ? ", msg.SenderName, msg.ReceiverName, msg.State, msg.Timestamp).Update("state", newState)
}
