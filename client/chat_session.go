package main

import (
	"fmt"
	"math/big"

	//"github.com/grigagod/chat-example/crypto"

	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
)

func (c *Client) StartChatSession() {
	c.friendInvites = make(map[string]*big.Int, 0)
	c.friends = make(map[string]*big.Int, 0)
	for {
		msg, err := c.wsReader.GetNext()
		if err != nil {
			fmt.Println(err)
			break
		}

		switch msg.Type {
		case websock.ChatInfoResponse:
			chatInfo := msg.Message.(*websock.ChatInfoMessage)
			c.users = chatInfo.Users
			fmt.Print(c.users)
		case websock.KeyExchangeStatus:
			keyExchInfo := msg.Message.(string)
			fmt.Print(keyExchInfo)
		case websock.KeyExchangeRequest:
			keyExchMsg := msg.Message.(*websock.KeyExchangeMessage)

			friendKey, err := util.UnmarshalKey(keyExchMsg.FriendPubKey)
			if err != nil {
				fmt.Println("Error while parsing other user pubKey")
			}
			c.friendInvites[keyExchMsg.Friendname] = friendKey

			fmt.Println("New friend request from ", keyExchMsg.Friendname)
		case websock.KeyExchangeAccept:
			friendData := msg.Message.(*websock.KeyExchangeMessage)

			friendKey, err := util.UnmarshalKey(friendData.FriendPubKey)
			if err != nil {
				fmt.Println("Error while parsing other user pubKey")
			}

			sharedKey := c.keys.KeyMixing(friendKey)
			c.friends[friendData.Friendname] = sharedKey
			delete(c.friendInvites, friendData.Friendname)

			fmt.Println("User:", friendData.Friendname, " accepted your invite")
		case websock.KeyExchangeDecline:
			friendName := msg.Message.(string)

			fmt.Println("User: ", friendName, "declined your invite")
		case websock.DirectMessage:
			message := msg.Message.(*websock.ChatMessage)

			if c.friends[message.Sender] != nil {

				decrMsg := util.DecryptDirectMessage(c.friends[message.Sender], message.Message)

				fmt.Println("[", message.Sender, " ", message.Timestamp, " : ", decrMsg)

			} else {
				fmt.Println("Can't decrypt entering message")
			}
		}
	}

}
