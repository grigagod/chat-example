package sdb

type Request struct {
	SenderName   string `gorm:"primaryKey"`
	SenderKey    []byte
	ReceiverName string `gorm:"primaryKey"`
	Receiver     User   `gorm:"foreignKey:ReceiverName;References:Username"`
}

func NewRequest(senderName string, publicKey []byte, receiverName string) *Request {
	return &Request{
		SenderName:   senderName,
		SenderKey:    publicKey,
		ReceiverName: receiverName,
	}
}
