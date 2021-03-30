package main

import (
	// "crypto/rsa"
	// "errors"
	// "io/ioutil"
	// "log"
	// "os"
	"bufio"
	"fmt"
	"github.com/grigagod/chat-example/crypto"
	"golang.org/x/net/websocket"
	"os"
	"strings"
)

type Client struct {
	ws      *websocket.Conn
	keys    crypto.Keys
	authKey []byte
}

func (c *Client) Connect(server string) bool {
	ws, err := websocket.Dial(server, "", "http://")
	if err != nil {
		fmt.Println("Error connecting to the server")
		return false
	}
	c.ws = ws

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
