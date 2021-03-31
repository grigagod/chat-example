package server

import (
	"github.com/grigagod/chat-example/pdb"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
)

func (s *Server) ResponseUsersInfo(ws *websocket.Conn) {
	var users []pdb.User

	s.Db.Select("Username").Find(&users)

	result := make(map[string]bool, len(users))

	for _, user := range users {
		result[user.Username] = true
	}

	response := &websock.ChatUsersMessage{
		Users: result,
	}

	websock.Send(ws, &websock.Message{Type: websock.ChatUsersResponse, Message: response})
}
