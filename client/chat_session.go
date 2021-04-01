package main

import (
	"fmt"

	"github.com/grigagod/chat-example/crypto"
	"github.com/grigagod/chat-example/pdb"
	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
)

func (c *Client) StartChatSession() {

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
			keyExchInfo := msg.Message.(pdb.NotificationState)
			fmt.Print("KeyExchangeStatusCode:", keyExchInfo)
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

			fmt.Println("User:", friendData.Friendname, " accepted your invite")
		case websock.KeyExchangeDecline:
			friendName := msg.Message.(string)
			fmt.Println("User: ", friendName, "declined your invite")
		}
	}

}
