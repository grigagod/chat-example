package server

import (
	"github.com/grigagod/chat-example/pdb"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
)

func (s *Server) ResponseUsersInfo(ws *websocket.Conn) {
	var users []pdb.User

	s.Db.Select("Username").Find(&users)

	result := make([]string, len(users))

	for i, user := range users {
		result[i] = user.Username
	}

	response := &websock.ChatInfoMessage{
		Users: result,
	}

	websock.Send(ws, &websock.Message{Type: websock.ChatInfoResponse, Message: response})
}

func (s *Server) ResponseKeyExchInit(ws *websocket.Conn, receivername string) {
	var sender pdb.User
	var receiver pdb.User

	if user, ok := s.Users.Get(ws); ok {
		s.Db.First(&sender, "username = ?", user.Username)

		if err := s.Db.First(&receiver, "username = ?", receivername).Error; err != nil {
			websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Failed to find user in db"})
		} else {
			notification := pdb.NewNotification(&sender, &receiver)

			if err := s.Db.Create(&notification).Error; err != nil {
				websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Failed to create invite in db"})
			} else {
				websock.Send(ws, &websock.Message{Type: websock.KeyExchangeStatus, Message: "Invite  is sent"})

				if ok := s.SendMessageIfActive(receivername, &websock.Message{Type: websock.KeyExchangeRequest, Message: &websock.KeyExchangeMessage{
					Friendname:   sender.Username,
					FriendPubKey: sender.PublicKey,
				}}); ok {
					go s.Db.Model(&notification).Update("state", pdb.Received)
				}
			}

		}
	}
}

func (s *Server) HandleKeyExchResponse(ws *websocket.Conn, msg *websock.Message) {
	var sender pdb.User
	var receiver pdb.User
	sendername := msg.Message.(string)

	user, _ := s.Users.Get(ws)

	s.Db.First(&receiver, "username = ?", user.Username)

	s.Db.First(&sender, "username = ?", sendername)

	var notification pdb.Notification

	s.Db.First(&notification, "sender_name = ? AND receiver_name =?", sender.Username, receiver.Username)
	switch msg.Type {
	case websock.KeyExchangeAccept:
		if ok := s.SendMessageIfActive(sendername, &websock.Message{Type: websock.KeyExchangeAccept, Message: &websock.KeyExchangeMessage{
			Friendname:   receiver.Username,
			FriendPubKey: receiver.PublicKey,
		}}); ok {
			go s.Db.Delete(&notification)

		} else {

			go s.Db.Model(&notification).Update("state", pdb.Accepted)
		}
	case websock.KeyExchangeDecline:
		if ok := s.SendMessageIfActive(sendername, &websock.Message{Type: websock.KeyExchangeDecline, Message: receiver.Username}); ok {
			go s.Db.Delete(&notification)

		} else {
			go s.Db.Model(&notification).Update("state", pdb.Declined)
		}
	}
}

func (s *Server) HandleDirectMessage(ws *websocket.Conn, msg *websock.Message) {

	chatMsg := msg.Message.(*websock.ChatMessage)
	var sender pdb.User
	var receiver pdb.User

	s.Db.First(&sender, "username = ?", chatMsg.Sender)

	s.Db.First(&receiver, "username = ?", chatMsg.Receiver)

	message := pdb.NewMessage(&sender, &receiver, chatMsg.Timestamp, chatMsg.Message)

	s.Db.Create(&message)

	if ok := s.SendMessageIfActive(chatMsg.Receiver, msg); ok {
		go s.Db.Model(&message).Update("state", pdb.MsgReceived)
	}
}
