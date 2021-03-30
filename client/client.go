package main

import (
	// "crypto/rsa"
	// "errors"
	// "io/ioutil"
	// "log"
	// "os"
	"bufio"
	"errors"
	"fmt"
	"github.com/grigagod/chat-example/crypto"
	"github.com/grigagod/chat-example/websock"
	"golang.org/x/net/websocket"
	"os"
	"strings"
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
	wsReader *WSReader
	ws       *websocket.Conn
	keys     crypto.Keys
	authKey  []byte
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

func menu() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "exit" {
			fmt.Println("You have been disconnected from the server")
			break
		}

		switch text {
		// TODO: make cases
		default:
			fmt.Println("Unknown command")
		}
	}
}

func main() {
	server := "ws://127.0.0.1:8001"
	var client Client
	if client.Connect(server) {
		menu()
	}

}
