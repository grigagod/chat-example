package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AddToFriendGUI struct {
	*GUI
	inviteFriendHandler func(username string)
	layout              *tview.Flex
	usersList           *tview.List
	users               []string

	focusableElements []tview.Primitive
	focusedIndex      int
}

func (gui *AddToFriendGUI) Create() {
	gui.usersList = tview.NewList()
	gui.usersList.
		SetSelectedFunc(gui.onUserSelected).
		SetBorder(true).
		SetTitle("Users")
	if len(gui.users) != 0 {
		for _, user := range gui.users {
			gui.usersList.AddItem(user, "", 0, nil)
		}
	}

	gui.layout = tview.NewFlex().SetFullScreen(true).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(gui.usersList, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

}
func (gui *AddToFriendGUI) onUserSelected(index int, name, secText string, scut rune) {
	gui.app.SetInputCapture(gui.chatGUI.KeyHandler)
	gui.pages.RemovePage("addFriend")
	gui.app.ForceDraw()
	gui.inviteFriendHandler(name)
}

func (gui *AddToFriendGUI) KeyHandler(key *tcell.EventKey) *tcell.EventKey {
	if key.Key() == tcell.KeyEsc {
		gui.pages.RemovePage("addFriend")
		gui.app.SetFocus(gui.pages)
		gui.app.SetInputCapture(gui.chatGUI.KeyHandler)
		gui.app.ForceDraw()
	}

	return key
}
