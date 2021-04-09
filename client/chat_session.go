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
	c.friendRequests = make(map[string]*big.Int, 0)
	c.friends = make(map[string]*big.Int, 0)

	for _, friend := range c.GetFriendsList() {
		c.friends[friend.FriendName] = new(big.Int).SetBytes(friend.SharedKey)
		c.gui.chatGUI.addToFriendList(friend.FriendName)
	}
	for _, request := range c.GetRequestsList() {
		c.friendRequests[request.SenderName] = new(big.Int).SetBytes(request.SenderKey)
		c.gui.chatGUI.addToRequestsList(request.SenderName)
	}

	for {
		msg, err := c.wsReader.GetNext()
		if err != nil {
			c.gui.ShowDialog((err.Error()), nil)
			continue
		}

		switch msg.Type {
		case websock.ChatInfoResponse:
			chatInfo := msg.Message.(*websock.ChatInfoMessage)
			filteredUsers := make([]string, 0)
			for _, user := range chatInfo.Users {
				if user == c.username {
					continue
				}
				i := 0
				for friend, _ := range c.friends {
					if user == friend {
						log.Println(friend)
						break
					} else {
						i = i + 1
					}

				}
				if len(c.friends) == i {
					filteredUsers = append(filteredUsers, user)
				}
			}

			c.users = filteredUsers
			c.gui.ShowAddFriendGUI(c)
		case websock.KeyExchangeStatus:
			keyExchInfo := msg.Message.(string)
			c.gui.ShowDialog(keyExchInfo, nil)
		case websock.KeyExchangeRequest:
			keyExchMsg := msg.Message.(*websock.KeyExchangeMessage)

			friendKey, err := util.UnmarshalKey(keyExchMsg.FriendPubKey)
			if err != nil {
				c.gui.ShowDialog("Error while parsing other user pubKey", nil)
				continue
			}
			c.friendRequests[keyExchMsg.Friendname] = friendKey

			c.gui.chatGUI.addToRequestsList(keyExchMsg.Friendname)
			go c.InsertIntoRequests(keyExchMsg.Friendname, keyExchMsg.FriendPubKey)
		case websock.KeyExchangeAccept:
			friendData := msg.Message.(*websock.KeyExchangeMessage)

			friendKey, err := util.UnmarshalKey(friendData.FriendPubKey)
			if err != nil {
				c.gui.ShowDialog("Error while parsing other user pubKey", nil)
				continue

			}

			sharedKey := c.keys.KeyMixing(friendKey)
			c.friends[friendData.Friendname] = sharedKey
			delete(c.friendRequests, friendData.Friendname)
			c.gui.chatGUI.addToFriendList(friendData.Friendname)

			go c.InsertIntoFriends(friendData.Friendname, sharedKey)
			go c.DeleteFromRequests(friendData.Friendname)
		case websock.KeyExchangeDecline:
			friendName := msg.Message.(string)
			c.gui.ShowDialog(fmt.Sprint("User :", friendName, " declined your invite"), nil)
		case websock.DirectMessage:
			message := msg.Message.(*websock.ChatMessage)

			if c.friends[message.Sender] != nil {
				decrMsg := util.DecryptDirectMessage(c.friends[message.Sender], message.Message)

				go c.dal.InsertIntoMessages(message.Sender, message.Receiver, decrMsg, message.Timestamp)
				if c.gui.chatGUI.selectedFriendName == message.Sender {

					c.gui.chatGUI.DisplayMessage(message.Sender, decrMsg, message.Timestamp)
				}

			} else {
				log.Println("Can't decrypt entering message")
			}
		}
	}
}
