package main

import (
	"github.com/gorilla/websocket"
	"github.com/grigagod/chat-example/pdb"
	"gorm.io/gorm"
	"log"
	"sync/atomic"
	"time"
)

type Config struct {
	DSN        string
	Keepaplive int
}

type Server struct {
	Config
	Db    *gorm.DB
	Users Users
}

func CreateServer(config Config) *Server {
	db, err := pdb.CreateConnection(config.DSN)

	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		Config: config,
		Db:     db,
		Users:  Users{data: make(map[*websocket.Conn]*User)},
	}
	return s
}

// AddClient adds a new client to Users
func (s *Server) AddClient(ws *websocket.Conn, user *User) {
	if !s.Users.Insert(ws, user) {
		log.Print("Websocket connection is already associated with a user")
	}
}

// RemoveClient removes a client from the ConnectedClients map
func (s *Server) RemoveClient(ws *websocket.Conn) {
	user, ok := s.Users.Remove(ws)
	if !ok {
		log.Print("Websocket was not in users-map")
		return
	}
	if user == nil {
		log.Print("Websocket was not associated with a user")
	}
}

// Pinger sends a ping message to the client in the interval specified in Keepalive in the ServerConfig
// If no pongs were received during the elapsed time, the server will close the client connection.
func (s *Server) Pinger(ws *websocket.Conn) (*time.Ticker, *int64) {
	ticker := time.NewTicker(time.Duration(s.Keepalive) * time.Second)
	pongCount := int64(1)

	go func() {
		for range ticker.C {
			if atomic.LoadInt64(&pongCount) == 0 {
				log.Printf("Client %s did not respond to ping in time", ws.Request().RemoteAddr)
				ws.Close()
				return
			}

			websock.Send(ws, &websock.Message{Type: websock.Ping})
			atomic.StoreInt64(&pongCount, 0)
		}
	}()

	return ticker, &pongCount
}
