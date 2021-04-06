package main

import (
	"github.com/rivo/tview"
	"math/big"
)

type GUIConfig struct {
	DefaultServerText   string
	createUserHandler   func(server, username string)
	loginUserHandler    func(server, username string)
	inviteFriendHandler func(friendname string)
	addToFriendsHandler func(server string, friendKey *big.Int)
}

type GUI struct {
	app      	 *tview.Application
	pages    	 *tview.Pages
	loginGUI 	 *LoginGUI
	chatGUI  	 *ChatGUI
	AddFriendGUI *AddToFriendGUI
}

func NewGUI(config *GUIConfig) *GUI {
	g := &GUI{
		app: tview.NewApplication(),
	}

	g.loginGUI = &LoginGUI{
		GUI:               g,
		DefaultServerText: config.DefaultServerText,
		createUserHandler: config.createUserHandler,
		loginUserHandler:  config.loginUserHandler}
	g.loginGUI.Create()

	g.chatGUI = &ChatGUI{
		GUI: 		g,
		addToFriendsHandler: config.addToFriendsHandler}
	g.chatGUI.Create()

	g.pages = tview.NewPages().
		AddPage("login", g.loginGUI.layout, true, true).
		AddPage("chat", g.chatGUI.layout, true, false).
		AddPage("addFriend", g.chatGUI.layout, true, false)

	g.app.SetRoot(g.pages, true).
		SetFocus(g.pages).
		SetInputCapture(g.loginGUI.KeyHandler)

	return g
}

func (g *GUI) ShowDialog(message string, onDismiss func()) {
	modal := tview.NewModal()
	modal.SetText(message).AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				g.pages.RemovePage("error")
			}
		})

	if onDismiss != nil {
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				onDismiss()
			}
		})
	}
	g.pages.AddPage("error", modal, true, true)
	g.app.SetFocus(modal)
}

func (g *GUI) ShowChatGUI(c *Client) {
	g.pages.SwitchToPage("chat")
	g.app.SetInputCapture(g.chatGUI.KeyHandler)
	g.ShowDialog("Welcome to chat", nil)
	go c.StartChatSession()
}

func (g *GUI) ShowAddFriendGUI(c *Client) {
	g.pages.SwitchToPage("addFriend")
	//g.app.SetInputCapture(g.AddFriendGUI.KeyHandler)
}

