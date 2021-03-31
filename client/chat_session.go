package main

import (
	"fmt"
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
		}
	}

}
