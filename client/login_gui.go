package main

// gui for log in page

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type LoginGUI struct {
	*GUI
	DefaultServerText string
	createUserHandler func(server string, username string)
	loginUserHandler  func(server string, username string)

	layout        *tview.Grid
	serverInput   *tview.InputField
	userNameInput *tview.InputField
	createBtn     *tview.Button
	loginBtn      *tview.Button
	statusText    *tview.TextView

	focusableElements []tview.Primitive
	focusIndex        int
}

func (gui *LoginGUI) Create() {
	gui.serverInput = tview.NewInputField().
		SetLabel("Server   ").
		SetFieldWidth(60).
		SetText(gui.DefaultServerText)

	gui.userNameInput = tview.NewInputField().
		SetLabel("Username ").
		SetFieldWidth(60)

	gui.createBtn = tview.NewButton("Create User")
	gui.loginBtn = tview.NewButton("Log In")

	gui.statusText = tview.NewTextView().
		SetTextColor(tcell.ColorLightBlue).
		SetTextAlign(tview.AlignCenter)

	gui.statusText.SetText("Welcome. Create a new user or log in using you private key file.")

	gui.layout = tview.NewGrid()
	gui.layout.SetRows(0, 1, 1, 1, 1, 0, 2).
		SetColumns(0, 30, 5, 30, 0).
		AddItem(gui.serverInput, 1, 1, 1, 3, 0, 0, false).
		AddItem(gui.userNameInput, 2, 1, 1, 3, 0, 0, true).
		AddItem(gui.createBtn, 4, 1, 1, 1, 0, 0, false).
		AddItem(gui.loginBtn, 4, 3, 1, 1, 0, 0, false).
		AddItem(gui.statusText, 6, 0, 1, 5, 0, 0, false).
		SetBorder(true).
		SetTitle("Chat Client")

	gui.focusableElements = []tview.Primitive{
		gui.serverInput, gui.userNameInput,
		gui.createBtn, gui.loginBtn}
	gui.focusIndex = 1
}

func (gui *LoginGUI) KeyHandler(ev *tcell.EventKey) *tcell.EventKey {
	if ev.Key() == tcell.KeyTab {
		gui.focusIndex++
		if gui.focusIndex == len(gui.focusableElements) {
			gui.focusIndex = 0
		}

		gui.app.SetFocus(gui.focusableElements[gui.focusIndex])
	} else if ev.Key == tcell.KeyEnter {
		switch gui.app.GetFocus() {
		case gui.createBtn:
			gui.createUserHandler(gui.serverInput.GetText(), gui.userNameInput.GetText())
		case gui.loginBtn:
			gui.loginUserHandler(gui.serverInput.GetText(), gui.userNameInput.GetText())
		}
	}

	return ev
}
