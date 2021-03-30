package main

import (
	"fmt"
	"github.com/grigagod/chat-example/crypto"
	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
	"io/ioutil"
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

// Called when the user pressed the "login user" button
// TODO: Refactor the huge function
func (c *Client) loginUserHandler(server string, username string) {
	if !c.Connect(server) {
		return
	}

	// Read private key from file
	message, err := ioutil.ReadFile(username + ".chat")
	if err != nil {
		fmt.Println("Error reading privatekey.pem file")
		return
	}

	privKey, err := util.UnmarshalKey(message)
	if err != nil {
		fmt.Println("Error parsing private key")
		return
	}

	// Send log in request to server
	websock.Send(c.ws, &websock.Message{Type: websock.LoginUser, Message: username})

	// Receive auth challenge from server
	res, err := c.wsReader.GetNext()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Got auth challenge")

	// Try to decrypt auth challenge
	keys := crypto.KeysFromPrivate(privKey)
	decKey := util.DecryptChallenge(keys.PublicKey, res.Message.([]byte))

	fmt.Println("DecryptChallenge in process")

	// Send decrypted auth key to server
	websock.Send(c.ws, &websock.Message{Type: websock.AuthChallengeResponse, Message: decKey})

	// Check response from server
	if res, err = c.wsReader.GetNext(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	// Login success, show the chat rooms GUI
	c.authKey = decKey
	c.keys = keys
	fmt.Println("Successfully logged in")
}
