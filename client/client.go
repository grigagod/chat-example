package main

import (
	// "crypto/rsa"
	// "errors"
	"io/ioutil"
	// "log"
	// "os"
	"bufio"
	"errors"
	"fmt"
	"github.com/grigagod/chat-example/crypto"
	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
	"math/big"
	"os"
	"strings"
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
	keys          *crypto.Keys
	authKey       []byte
	username      string
	users         []string
	friends       map[string]*big.Int
	friendInvites map[string]*big.Int
}

func (c *Client) Connect(server string) bool {
	if c.wsReader == nil {
		ws, err := websocket.Dial(server, "", "http://")
		if err != nil {
			fmt.Println("Error connecting to server")
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

func savePrivKey(username string, privKey *big.Int) {
	message := util.MarshalKey(privKey)
	if err := ioutil.WriteFile(username+".chat", message, 0644); err != nil {
		fmt.Println(err)
	}
}
func menu(c *Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		args := strings.Split(text, " ")
		cmd := strings.TrimSpace(args[0])

		if cmd == "/exit" {
			fmt.Println("You have been disconnected from the server")
			break
		}

		switch cmd {
		case "/name":
			c.username = args[1]
			fmt.Println("Your username is : ", c.username)
		case "/register":
			c.createUserHandler(serverStr, c.username)
		case "/login":
			if args[1] != "" {
				c.loginUserHandler(serverStr, args[1])

			}
		case "/cinfo":
			c.chatInfoHandler()
		case "/invite":
			for _, v := range c.users {
				if v == args[1] {
					c.addToFriendsHandler(args[1])
				}
			}
			fmt.Println("No registered user with such nickname")
		case "/accept":
			for k := range c.friends {
				if k == args[1] {

				}
			}
		default:
			fmt.Println("Unknown command")
		}
	}
}

func main() {
	client := &Client{}

	if client.Connect(serverStr) {
		menu(client)
	}

}
