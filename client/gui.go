package main

import (
  "time"

  "github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type GUIConfig struct {
  DefaultServerText string
  createUserHandler func(server string, username string)
  loginUserHandler func(server string, username string)
  addToFriendsHandler func(friendname string)
}

type GUI struct {
  app *tview.Application
  pages *tview.Pages
  loginGUI *LoginGUI
  chatGUI *ChatGUI
}

func NewGUI(config *GUIConfig) *GUI {
  g := &GUI{
    app: tview.NewApplication()
  }

  g.loginGUI = &LoginGUI{
    GUI: g,
    createUserHandler: config.createUserHandler,
    loginUserHandler: config.loginUserHandler,
  }
  g.loginGUI.Create()

  g.ChatGUI = &ChatGUI {GUI: g}
  g.ChatGUI.Create()

  g.pages = tview.NewPages().AddPage("login", g.LoginGUI.layout, true, true).
  AddPage("chat", g.ChatGUI.layout, true, false)

  g.app.SetRoot(g.pages, true).SetFocus(g.pages).SetInputCapture(g.loginGUI.KeyHandler)

  return g
}
