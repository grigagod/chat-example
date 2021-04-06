package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"log"
)

type AddToFriendGUI struct {
	*GUI
	inviteFriendHandler func(username string)
	layout 	  *tview.Grid
	usersList *tview.List
	addBtn 	  *tview.Button
	users	  []string

	focusableElements 		 []tview.Primitive
	focusedIndex      		 int
}


func (gui *AddToFriendGUI) Create() {
	gui.usersList = tview.NewList()
	gui.usersList.
		SetSelectedFunc(gui.onUserSelected).
		SetBorder(true).
		SetTitle("Users")
	for _, user := range gui.users {
		gui.usersList.AddItem(user, "", 0, nil)
	}


	gui.addBtn = tview.NewButton("Add To Friends")
	gui.layout = tview.NewGrid()
	gui.layout.
		SetRows(0, 3, 1).
		SetColumns(20, 1, 20, 0, 30).
		SetBorders(true).
		AddItem(gui.usersList, 0, 4, 2, 1, 0, 0, false).
		AddItem(gui.addBtn, 2, 0, 1, 1, 0, 0, false)
	
	gui.focusableElements = []tview.Primitive{
		gui.usersList,
		gui.addBtn}
	gui.focusedIndex = 1

}

func (gui *AddToFriendGUI) onUserSelected(index int, name, secText string, scut rune) {
	gui.CurrentChatName = name
	log.Println(gui.CurrentChatName)
	gui.inviteFriendHandler(name)	
}


func (gui *AddToFriendGUI) KeyHandler(key *tcell.EventKey) *tcell.EventKey {
	if key.Key() == tcell.KeyEsc {
		gui.pages.SwitchToPage("chat")
		gui.app.SetInputCapture(gui.chatGUI.KeyHandler)
	}
	if key.Key() == tcell.KeyTab {
		gui.focusedIndex++
		if gui.focusedIndex == len(gui.focusableElements) {
			gui.focusedIndex = 0
		}
		gui.app.SetFocus(gui.focusableElements[gui.focusedIndex])
	} else if key.Key() == tcell.KeyEnter {
		switch gui.app.GetFocus() {
		case gui.addBtn:
			return nil
		case gui.usersList:
			return nil
		}
	}
	
	return key
}
