package main

import (
	// "crypto/rsa"
	// "errors"
	"log"

	// "log"
	// "os"

	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/grigagod/chat-example/crypto"

	//"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
)

const (
	serverStr = "ws://127.0.0.1:8001"
)

// Result is used by WSReader to communicate the websocket messages between threads
type Result struct {
	Message *websock.Message
	Err     error
}

// WSReader reads messages from the websocket in the background
type WSReader struct {
	OnDisconnect func()
	Ws           *websocket.Conn
	c            chan Result
}

// Reader runs in a separate goroutine and listens for messages on the websocket
func (wr *WSReader) Reader() {
	for {
		msg := new(websock.Message)
		err := websock.Receive(wr.Ws, msg)
		if err != nil {
			fmt.Println(err)
			break
		}

		if msg.Type == websock.Ping {
			websock.Send(wr.Ws, &websock.Message{Type: websock.Pong})
		} else if msg.Type == websock.Error {
			wr.c <- Result{Message: nil, Err: errors.New(msg.Message.(string))}
		} else {
			wr.c <- Result{Message: msg, Err: nil}
		}
	}

	wr.OnDisconnect()
}

// GetNext retrieves the next websocket message from the message pool
func (wr *WSReader) GetNext() (*websock.Message, error) {
	result := <-wr.c
	return result.Message, result.Err
}

type Client struct {
	wsReader      *WSReader
	ws            *websocket.Conn
	dal           *DAL
	keys          *crypto.Keys
	authKey       []byte
	username      string
	users         []string
	friends       map[string]*big.Int
	friendInvites map[string]*big.Int
	gui           *GUI
}

func (c *Client) Connect(server string) bool {
	if c.wsReader == nil {
		ws, err := websocket.Dial(server, "", "http://")
		if err != nil {
			c.gui.ShowDialog("Error connecting to server", nil)
			return false
		}

		c.ws = ws
		c.wsReader = &WSReader{
			OnDisconnect: func() {
				fmt.Println("Disconnected from the server")
			},
			Ws: ws,
			c:  make(chan Result, 10)}
		go c.wsReader.Reader()
	}
	return true
}

func main() {
	f, _ := os.OpenFile("client_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)

	client := &Client{}
	client.dal = createDAL("chat.db")

	guiConfig := &GUIConfig{
		DefaultServerText:   serverStr,
		createUserHandler:   client.createUserHandler,
		loginUserHandler:    client.loginUserHandler,
		inviteFriendHandler: client.inviteFriendHandler,
	}

	client.gui = NewGUI(guiConfig)

	// Enter GUI event loop
	if err := client.gui.app.Run(); err != nil {
		log.Fatal(err)
	}
}
