package websock

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"

	"golang.org/x/net/websocket"
)

// codec is the decoder for messages sent over the websocket
var codec = websocket.Codec{Marshal: marshalMessage, Unmarshal: unmarshalMessage}

// Send sends a message to the connection synchronously
//
// NB! This will fail if the given message contains incompatible type and content
func Send(ws *websocket.Conn, msg *Message) error {
	return codec.Send(ws, msg)
}

// Receive fetches a message from the connection synchronously
//
// NB! This will fail if the received message contains incompatible type and content
func Receive(ws *websocket.Conn, msg *Message) error {
	return codec.Receive(ws, msg)
}

// Register types for gob encoding/decoding
func init() {
	gob.Register(&RegisterUserMessage{})
	gob.Register(&ChatInfoMessage{})
	gob.Register(&KeyExchangeMessage{})
	gob.Register(&ChatMessage{})
}

func marshalMessage(v interface{}) ([]byte, byte, error) {
	msg, ok := v.(*Message)
	if !ok {
		return nil, websocket.TextFrame, errors.New("Input to marshalMessage was not of type *Message")
	}

	if err := checkType(msg.Message, msg.Type); err != nil {
		return nil, websocket.TextFrame, err
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(msg); err != nil {
		return nil, websocket.TextFrame, err
	}

	return buf.Bytes(), websocket.TextFrame, nil
}

func unmarshalMessage(data []byte, payloadType byte, v interface{}) error {
	msg, ok := v.(*Message)
	if !ok {
		return errors.New("Input to unmarshalMessage was not of type *Message")
	}

	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(msg); err != nil {
		return err
	}

	if err := checkType(msg.Message, msg.Type); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func checkType(v interface{}, msgType MessageType) error {
	switch msgType {
	case Error, OK, LoginUser, KeyExchangeInit, KeyExchangeStatus, KeyExchangeDecline:
		if _, ok := v.(string); !ok {
			return errors.New("Expected message type string")
		}

	case RegisterUser:
		if _, ok := v.(*RegisterUserMessage); !ok {
			return errors.New("Expected message type *RegisterUserMessage")
		}

	case AuthChallenge, AuthChallengeResponse:
		if _, ok := v.([]byte); !ok {
			return errors.New("Expected message type []byte")
		}

	case DirectMessage:
		if _, ok := v.(*ChatMessage); !ok {
			return errors.New("Expected message type *ChatMessage")
		}
	case KeyExchangeAccept:
		_, ok1 := v.(string)
		if _, ok2 := v.(*KeyExchangeMessage); !ok1 && !ok2 {
			return errors.New("Expected message type string or *KeyExchangeMessage")
		}
	case KeyExchangeRequest:
		if _, ok := v.(*KeyExchangeMessage); !ok {
			return errors.New("Expected message type *KeyExchangeMessage")
		}
	case ChatUsersInfo, Ping, Pong:
		if v != nil {
			return errors.New("Expected message to be nil")
		}
	case ChatInfoResponse:
		if _, ok := v.(*ChatInfoMessage); !ok {
			return errors.New("Expected message type *ChatUsersMessage")
		}
	default:
		return errors.New("Invalid message type")
	}

	return nil
}
