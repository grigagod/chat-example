package sdb

type MessageState int

const (
	MsgSent MessageState = iota

	MsgReceived

	MsgFailed
)

type Message struct {
	SenderName   string       `gorm:"index"`
	ReceiverName string       `gorm:"index"`
	State        MessageState `gorm:"index"`
	Message      string
	Timestamp    int64
}

func newSentMessage(senderName string, receiverName string, message string, timestamp int64) *Message {
	return &Message{
		SenderName:   senderName,
		ReceiverName: receiverName,
		State:        MsgSent,
		Message:      message,
		Timestamp:    timestamp,
	}
}

func newReceivedMessage(senderName string, receiverName string, message string, timestamp int64) *Message {
	return &Message{
		SenderName:   senderName,
		ReceiverName: receiverName,
		State:        MsgReceived,
		Message:      message,
		Timestamp:    timestamp,
	}
}
