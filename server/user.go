package server

import (
	"bytes"
	"errors"
	"log"
	"math/big"
	"sync"

	"github.com/grigagod/chat-example/pdb"
	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
	"gorm.io/gorm"
)

type Users struct {
	sync.Mutex

	data map[*websocket.Conn]*User
}

// Get gets the User of a connected websocket client
//
// Returns true on success and false on missing user
func (users *Users) Get(ws *websocket.Conn) (user *User, ok bool) {
	users.Lock()
	defer users.Unlock()
	user, ok = users.data[ws]
	return
}

// GetByName gets the User of a connected websocket client by username
//
// Returns true on success and false on missing user
func (users *Users) GetWSByName(username string) (ws *websocket.Conn, ok bool) {
	users.Lock()
	defer users.Unlock()

	for ws, user := range users.data {
		if user.Username == username {
			return ws, true
		}
	}
	return nil, false
}

// Remove deletes the connection between a websocket and a user
func (users *Users) Remove(ws *websocket.Conn) (user *User, ok bool) {
	users.Lock()
	defer users.Unlock()

	user, ok = users.data[ws]
	if ok {
		delete(users.data, ws)
	}
	return
}

// Insert adds the given User to the collection indexed by the websocket
// connection
//
// Returns true on success and false on already existing association between
// socket and user
func (users *Users) Insert(ws *websocket.Conn, user *User) bool {
	users.Lock()
	defer users.Unlock()

	// This connection already has an associated user
	if user, ok := users.data[ws]; ok && user != nil {
		return false
	}
	users.data[ws] = user
	return true
}

// ForEach performs the given function on all stored users
func (users *Users) ForEach(f func(*websocket.Conn, *User)) {
	users.Lock()
	defer users.Unlock()
	for ws, user := range users.data {
		f(ws, user)
	}
}

// Len gets the amount of registered users
func (users *Users) Len() int {
	users.Lock()
	defer users.Unlock()
	return len(users.data)
}

// KeyMatches checks that an authentication key matches the one for this user
func (u *User) KeyMatches(authKey []byte) bool {
	return bytes.Equal(u.AuthKey, authKey)
}

type User struct {
	sync.Mutex
	Username  string
	AuthKey   []byte
	PublicKey *big.Int
}

// RegisterUser registers a new user, and adds it to the database
func (s *Server) RegisterUser(ws *websocket.Conn, msg *websock.RegisterUserMessage) {
	// Add new user to database
	user := pdb.NewUser(msg.Username, msg.PublicKey)
	if err := s.Db.Create(&user).Error; err != nil {
		websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Error registering user"})
		return
	}

	websock.Send(ws, &websock.Message{Type: websock.OK, Message: "User registered"})
}

// LoginUser authenticates a user using a randomly generated authentication token
// This token is encrypted with the public key of the username the client is trying to log in as
// The client is then expected to respond with the correct decrypted token
// TODO check if user is already logged in
func (s *Server) LoginUser(ws *websocket.Conn, username string) bool {
	// Create new user object
	newUser, encKey, err := NewUser(s.Db, username)
	if err != nil {
		websock.Send(ws, &websock.Message{Type: websock.Error, Message: "User does not exist"})
		return false
	}

	// Send auth challenge
	websock.Send(ws, &websock.Message{Type: websock.AuthChallenge, Message: encKey})

	// Receive auth challenge response
	res := new(websock.Message)
	if err := websock.Receive(ws, res); err != nil {
		log.Println(err)
		return false
	}

	// Check that the received decrypted key matches the original auth key
	if newUser.KeyMatches(res.Message.([]byte)) {
		log.Printf("Client %s authenticated as user %s\n", ws.Request().RemoteAddr, newUser.Username)
		s.AddClient(ws, newUser)
		websock.Send(ws, &websock.Message{Type: websock.OK, Message: "Logged in"})
		return true
	} else {
		log.Printf("Client %s failed authentication", ws.Request().RemoteAddr)
	}

	websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Invalid auth key"})
	return false
}

// NewUser creates a new user object for a connected client, with the username, generated (temporary) authentication
// key and the encrypted version of the key. A random byte slice is generated and encrypted with the users public key, the user
// is expected to send in response the decrypted string
func NewUser(db *gorm.DB, username string) (*User, []byte, error) {
	// Retrieve user from DB

	user := new(pdb.User)
	if err := db.Where("Username = ?", username).First(&user).Error; err != nil {
		log.Println("No")
		return nil, nil, errors.New("Registered user not found")
	}

	// Unmarshal public key
	pubKey, err := util.UnmarshalKey(user.PublicKey)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	// Generate auth challenge
	encKey, authKey := util.GenAuthChallenge(pubKey)

	return &User{
		Username:  username,
		AuthKey:   authKey,
		PublicKey: pubKey}, encKey, nil
}

func (s *Server) SendExistingNotifications(ws *websocket.Conn) {
	user, _ := s.Users.Get(ws)
	var invites []pdb.Notification
	var receiver pdb.User

	// Finding invites which is sent by this user
	s.Db.Where("sender_name = ? AND  state IN ?", user.Username, []pdb.NotificationState{pdb.Accepted, pdb.Declined}).Find(&invites)

	for _, invite := range invites {
		switch invite.State {
		case pdb.Accepted:
			s.Db.First(&receiver, "username = ?", invite.ReceiverName)
			err := websock.Send(ws, &websock.Message{Type: websock.KeyExchangeAccept, Message: &websock.KeyExchangeMessage{
				Friendname:   receiver.Username,
				FriendPubKey: receiver.PublicKey,
			}})
			if err == nil {
				s.Db.Delete(&invite)
			}
		case pdb.Declined:
			err := websock.Send(ws, &websock.Message{Type: websock.KeyExchangeDecline, Message: invite.ReceiverName})

			if err == nil {
				s.Db.Delete(&invite)
			}
		}
	}

	var sender pdb.User
	var invitations []pdb.Notification
	// Finding invite which is sent to this user

	s.Db.Where("receiver_name = ? AND state = ? ", user.Username, pdb.Initiated).Find(&invitations)

	for _, invitation := range invitations {
		s.Db.First(&sender, "username = ?", invitation.SenderName)

		err := websock.Send(ws, &websock.Message{Type: websock.KeyExchangeRequest, Message: &websock.KeyExchangeMessage{
			Friendname:   sender.Username,
			FriendPubKey: sender.PublicKey,
		}})

		if err == nil {
			s.Db.Model(&invitation).Where("state = ?", pdb.Initiated).Update("state", pdb.Received)
		} else {
			log.Println(err)
		}
	}

	var messages []pdb.Message

	s.Db.Where("receiver_name =? AND state =?", user.Username, pdb.MsgInitiated).Find(&messages)

	for _, msg := range messages {
		err := websock.Send(ws, &websock.Message{Type: websock.DirectMessage, Message: toChatMessage(&msg)})
		if err == nil {
			s.Db.Model(&msg).Where("state = ?", pdb.MsgInitiated).Update("state", pdb.MsgReceived)
		}
	}

}

func (s *Server) SendMessageIfActive(username string, msg *websock.Message) (ok bool) {
	if ws, ok := s.Users.GetWSByName(username); !ok {
		ok = false
	} else {
		err := websock.Send(ws, msg)
		if err != nil {
			ok = false
		}

		ok = true
	}
	return ok
}

func toChatMessage(msg *pdb.Message) *websock.ChatMessage {
	return &websock.ChatMessage{
		Sender:    msg.SenderName,
		Receiver:  msg.ReceiverName,
		Timestamp: msg.Timestamp,
		Message:   msg.Message,
	}
}
