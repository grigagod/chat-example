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

	// KeyExchangeInit is sent to the sever when client wants to start KeyExchange with another client
	KeyExchangeInit
	// PubKeyResponse is sent by the server to the client in response to an KeyExchangeRequest
	KeyExchangeStatus

	// KeyExchangeChallenge is initiated by PubKeyRequest msg and sent by server to another in order to start Direct Messaging ("starting friendship")
	KeyExchangeRequest
	// KeyExchangeAccept is sent by the invited client in response to KeyExchangeRequest message
	KeyExchangeAccept

	//KeyExchangeDecline
	KeyExchangeDecline

	// SendDirect is sent when a client sends a direct message
	SendDirect
	// DirectMessageReceived is sent by the server when another user receive a direct message
	DirectMessageReceived

	// ChatUsersInfo in sent by the client to the server when client wants get list of all users in chat with their status
	ChatUsersInfo
	// ChatUsersResponse is sent by server to client in response to ChatUsersInfo message
	ChatInfoResponse

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

// ChatInfoMessage is sent by the server as a response for ChatUsersInfo
type ChatInfoMessage struct {
	Users []string
}

// KeyExchangeMessage is sent by the server
type KeyExchangeMessage struct {
	Friendname   string
	FriendPubKey []byte
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
