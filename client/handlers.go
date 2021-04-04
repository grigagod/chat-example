package main

import (
	"fmt"
	"io/ioutil"
	"math/big"

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
	c.username = username
	fmt.Println("Successfully logged in, starting session")
	go c.StartChatSession()
}

func (c *Client) chatInfoHandler() {
	websock.Send(c.ws, &websock.Message{Type: websock.ChatUsersInfo})

}

func (c *Client) addToFriendsHandler(friendname string, friendkey *big.Int) {
	err := websock.Send(c.ws, &websock.Message{Type: websock.KeyExchangeAccept, Message: friendname})
	if err != nil {
		fmt.Println(err)
	} else {
		sharedKey := c.keys.KeyMixing(friendkey)
		c.friends[friendname] = sharedKey

		delete(c.friendInvites, friendname)
	}

}

func (c *Client) inviteFriendHandler(friendname string) {
	err := websock.Send(c.ws, &websock.Message{Type: websock.KeyExchangeInit, Message: friendname})
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) declineFriendHandler(friendname string) {
	err := websock.Send(c.ws, &websock.Message{Type: websock.KeyExchangeDecline, Message: friendname})
	if err != nil {
		fmt.Println(err)
	} else {

		delete(c.friendInvites, friendname)
	}
}

func (c *Client) sendDirectMessage(friendname string, msg string) {
	timestamp := util.NowMillis()
	err := websock.Send(c.ws, &websock.Message{Type: websock.DirectMessage, Message: &websock.ChatMessage{
		Sender:    c.username,
		Receiver:  friendname,
		Timestamp: timestamp,
		Message:   util.EncryptDirectMessage(c.friends[friendname], msg),
	}})
	if err != nil {
		fmt.Println("Message is now sent")
	} else {
		fmt.Println("Message is sent")
	}
}
