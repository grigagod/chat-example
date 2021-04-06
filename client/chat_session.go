package main

import (
	"fmt"
	"log"
	"math/big"

	//"github.com/grigagod/chat-example/crypto"

	"github.com/grigagod/chat-example/util"
	"github.com/grigagod/chat-example/websock"
)

func (c *Client) StartChatSession() {
	c.friendInvites = make(map[string]*big.Int, 0)
	c.friends = make(map[string]*big.Int, 0)

	for _, friend := range c.dal.GetFriendsList(c.username) {
		c.friends[friend.FriendName] = new(big.Int).SetBytes(friend.SharedKey)
	}
	for _, request := range c.dal.GetRequestsList(c.username) {
		c.friendInvites[request.SenderName] = new(big.Int).SetBytes(request.SenderKey)
	}

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
			for idx, user := range c.users {
				if user == c.username {
					c.users = append(c.users[:idx], c.users[idx+1:]...)
				}
				for friend, _ := range c.friends {
					if user == friend {
						c.users = append(c.users[:idx], c.users[idx+1:]...)
					}
				}
			}
			//log.Println(c.users)
			c.gui.ShowAddFriendGUI(c)
		case websock.KeyExchangeStatus:
			keyExchInfo := msg.Message.(string)
			c.gui.ShowDialog(keyExchInfo, nil)
		case websock.KeyExchangeRequest:
			keyExchMsg := msg.Message.(*websock.KeyExchangeMessage)

			friendKey, err := util.UnmarshalKey(keyExchMsg.FriendPubKey)
			if err != nil {
				log.Println("Error while parsing other user pubKey")
			}
			c.friendInvites[keyExchMsg.Friendname] = friendKey

			log.Println("New friend request from ", keyExchMsg.Friendname)
			go c.dal.InsertIntoRequests(keyExchMsg.Friendname, keyExchMsg.FriendPubKey, c.username)
		case websock.KeyExchangeAccept:
			friendData := msg.Message.(*websock.KeyExchangeMessage)

			friendKey, err := util.UnmarshalKey(friendData.FriendPubKey)
			if err != nil {
				log.Println("Error while parsing other user pubKey")
			}

			sharedKey := c.keys.KeyMixing(friendKey)
			c.friends[friendData.Friendname] = sharedKey
			delete(c.friendInvites, friendData.Friendname)

			log.Println("User:", friendData.Friendname, " accepted your invite")

			go c.dal.InsertIntoFriends(friendData.Friendname, sharedKey, c.username)
			go c.dal.DeleteFromRequests(friendData.Friendname, c.username)
		case websock.KeyExchangeDecline:
			friendName := msg.Message.(string)

			log.Println("User: ", friendName, "declined your invite")
		case websock.DirectMessage:
			message := msg.Message.(*websock.ChatMessage)

			if c.friends[message.Sender] != nil {

				decrMsg := util.DecryptDirectMessage(c.friends[message.Sender], message.Message)

				log.Println("[", message.Sender, "]: ", decrMsg)
				go c.dal.InsertIntoMessages(message.Sender, message.Receiver, decrMsg, message.Timestamp)
			} else {
				log.Println("Can't decrypt entering message")
			}
		}
	}

}
