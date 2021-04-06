package main

import (
	"github.com/rivo/tview"
)

type AddToFriendGUI struct {
	*GUI
	addToFriendsHandler func(server string, username string)
	layout 	  *tview.Grid
	usersList *tview.List
	addBtn 	  *tview.Button
}


func (gui *AddToFriendGUI) Create() {
	
}