package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	//"math/big"
)

type AddToFriendGUI struct {
	*GUI
	inviteFriendHandler func(username string)
	layout 	  *tview.Flex
	usersList *tview.List
	users	  []string
}


func (gui *AddToFriendGUI) Create() {
	gui.usersList = tview.NewList()
	for _, user := range gui.users {
		gui.usersList.AddItem(user, "", 0, nil)
	}
	gui.usersList.
		SetBorder(true).
		SetTitle("Users")


	gui.layout = tview.NewFlex()
	gui.layout.
		AddItem(gui.usersList, 0, 1, false)
}

func (gui *AddToFriendGUI) KeyHandler(ev *tcell.EventKey) *tcell.EventKey {
	return ev
}
