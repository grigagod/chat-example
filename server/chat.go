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
