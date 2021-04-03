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
				websock.Send(ws, &websock.Message{Type: websock.KeyExchangeStatus, Message: pdb.Initiated})
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

	switch msg.Type {
	case websock.KeyExchangeAccept:
		s.Db.Model(&pdb.Notification{}).Where("sender_name = ?", sender.Username).Where("receiver_name = ?", receiver.Username).Update("state", pdb.Accepted)
	case websock.KeyExchangeDecline:
		s.Db.Model(&pdb.Notification{}).Where("sender_name = ?", sender.Username).Where("receiver_name = ?", receiver.Username).Update("state", pdb.Declined)
	}
}

func (s *Server) CheckForExistingInvites() {

}
