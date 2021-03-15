package main

import (
	"github.com/gorilla/websocket"
	"github.com/monnand/dhkx"
	"sync"
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

// KeyMatches checks that an authentication key matches the one for this user
func (u *User) KeyMatches(publicKey *dhkx.DHKey) bool {
	return u.PublicKey == publicKey
}

type Result struct {
	Message *websocket.PreparedMessage
}

type User struct {
	sync.Mutex
	Username  string
	PublicKey *dhkx.DHKey
}
