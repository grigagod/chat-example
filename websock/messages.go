package websock

type MessageType int

const (
	// Error means that en error occurred
	Error MessageType = iota
	// OK means that the action was successful
	OK

	// RegisterUser is sent when a client wants to register a new user
	RegisterUser
	// LoginUser is sent when a client wants to authenticate as a user
	LoginUser
	// AuthChallenge is sent by the server when an authentication challenge is initiated
	AuthChallenge
	// AuthChallengeResponse is sent by the client in resposne to an authentication challenge
	AuthChallengeResponse

	// PubKeyRequest is sent to the sever when client wants get another client public key
	PubKeyRequest
	// PubKeyResponse is sent by the server to the client in response to an KeyExchangeRequest
	PubKeyResponse

	// KeyExchangeChallenge is sent from one client to another in order to start Direct Messaging
	KeyExchangeChallenge
	// KeyExchangeResponse is sent by the another client in response to KeyExchangeChallenge
	KeyExchangeResponse

	// SendDirect is sent when a client sends a direct message
	SendDirect
	// DirectMessageReceived is sent by the server when another user receive a direct message
	DirectMessageReceived

	// UserJoined is sent by the server when a user joins a chat server
	UserJoined
	// UserLeft is sent by the server when a user leaves a chat server
	UserLeft

	// Ping is a keepalive message sent by the server
	Ping
	// Pong is sent by the client in response to a Ping message
	Pong
)

// Message is the "base" message which is used for all websocket messages
// Type contains the type of the message (one of the MessageType enums)
// Message contains the actual content of the message, which can be a string, byte slice, a struct, or nil.
type Message struct {
	Type    MessageType
	Message interface{}
}

// RegisterUserMessage is the message sent by a client to request user registration
type RegisterUserMessage struct {
	Username  string
	PublicKey []byte
}

// User is used in PublicKey , and by the server when notifying a client about a new connected user
type User struct {
	Username  string
	PublicKey []byte
}

// ChatMessage is used in DirectMessage, and by the server when notifying a client about a new chat message
type ChatMessage struct {
	Sender    string
	Timestamp int64
	Message   []byte
}

type SendChatMessage struct {
	EncryptedContent map[string][]byte
}
