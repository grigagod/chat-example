package main

import (
	"fmt"

	"github.com/grigagod/chat-example/crypto"
	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
)

// Called when user pressed the "create user" button
func (c *Client) createUserHandler(server string, username string) {
	if !c.Connect(server) {
		return
	}

	// Generate new key pair
	keys := new(crypto.Keys)
	keys.GenerateKeys()

	// Send a request to register the user
	regUserMsg := &websock.RegisterUserMessage{
		Username:  username,
		PublicKey: util.MarshalKey(keys.PublicKey)}

	websock.Send(c.ws, &websock.Message{Type: websock.RegisterUser, Message: regUserMsg})

	_, err := c.wsReader.GetNext()
	if err != nil {
		fmt.Println("Did not get a response from the server")
		return
	}

	// Save private key to file
	savePrivKey(username, keys.PrivateKey)

	fmt.Println("User created. You can now log in.")
}
