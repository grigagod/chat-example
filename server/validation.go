package server

import (
	"strings"
	"unicode"

	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
)

// Checks that a string only contains alphanumeric characters
func isAlphaNumeric(input string) bool {
	for _, ch := range input {
		if !unicode.IsLetter(ch) &&
			!unicode.IsNumber(ch) {
			return false
		}
	}
	return true
}

// ValidateRegisterUser validates the contents of a request from a client to
// register a new user. The length of the username and the public key bit-length is validated.
func ValidateRegisterUser(ws *websocket.Conn, msg *websock.RegisterUserMessage) bool {
	msg.Username = strings.TrimSpace(msg.Username)

	if len(msg.Username) < 3 {
		websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Username must contain at least 3 characters"})
		return false
	} else if len(msg.Username) > 20 {
		websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Username cannot contain more than 20 characters"})
		return false
	} else if !isAlphaNumeric(msg.Username) {
		websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Username can only contain alphanumeric characters"})
		return false
	}

	// Check key length
	if _, err := util.UnmarshalKey(msg.PublicKey); err != nil {
		websock.Send(ws, &websock.Message{Type: websock.Error, Message: "Invalid public key"})
		return false
	}
	return true
}
