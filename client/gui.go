package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GUIConfig struct {
	DefaultServerText   string
	createUserHandler   func(server string, username string)
	loginUserHandler    func(server string, username string)
	inviteFriendHandler func(friendname string)
}

type GUI struct {
	app      *tview.Application
	pages    *tview.Pages
	loginGUI *LoginGUI
	chatGUI  *ChatGUI
}

func NewGUI(config *GUIConfig) *GUI {
	g := &GUI{
		app: tview.NewApplication(),
	}

	g.loginGUI = &LoginGUI{
		GUI:               g,
		DefaultServerText: config.DefaultServerText,
		createUserHandler: config.createUserHandler,
		loginUserHandler:  config.loginUserHandler,
	}
	g.loginGUI.Create()

	g.chatGUI = &ChatGUI{GUI: g}
	g.chatGUI.Create()

	g.pages = tview.NewPages().AddPage("login", g.loginGUI.layout, true, true).
		AddPage("chat", g.chatGUI.layout, true, false)

	g.app.SetRoot(g.pages, true).SetFocus(g.pages).SetInputCapture(g.loginGUI.KeyHandler)

	return g
}

func (g *GUI) ShowDialog(message string, onDismiss func()) {
	modal := tview.NewModal()
	modal.SetText(message).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.pages.RemovePage("error")
		}).
		SetBackgroundColor(tcell.ColorDarkRed)

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
