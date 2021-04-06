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
	addBtn 	  *tview.Button
}


func (gui *AddToFriendGUI) Create() {
	gui.addBtn = tview.NewButton("Add")	
	gui.usersList = tview.NewList()
	gui.usersList.
		AddItem("asd", "", 0, nil).
		SetBorder(true).
		SetTitle("Users")

	gui.layout = tview.NewFlex()
	gui.layout.
		AddItem(gui.usersList, 0, 1, false)
}

func (gui *AddToFriendGUI) KeyHandler(ev *tcell.EventKey) *tcell.EventKey {
	return ev
}
