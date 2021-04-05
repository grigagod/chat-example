package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ChatGUI contains the widgets/state for the chat main window view
type ChatGUI struct {
	*GUI
	SendDirectMessageHandler func(friendname string, msg string)
	CurrentChatName          string
	LeaveChatHandler         func()
	layout                   *tview.Grid
	friendList               *tview.List
	msgView                  *tview.TextView
	msgInput                 *tview.InputField
	addFriendBtn             *tview.Button
	checkInvitesBtn          *tview.Button
}

// Create initializes the widgets in the chat GUI
func (gui *ChatGUI) Create() {
	gui.friendList = tview.NewList()
	gui.friendList.SetSelectedFunc(gui.onFriendSelected).
		SetTitle("Friends").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)

	gui.msgView = tview.NewTextView()
	gui.msgView.SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Chat")

	sendBtn := tview.NewButton("(Enter) Send")
	exitBtn := tview.NewButton("(Esc) Leave")

	gui.layout = tview.NewGrid()
	gui.layout.SetRows(0, 3, 1).
		SetColumns(20, 1, 20, 0, 30).
		AddItem(gui.msgView, 0, 0, 1, 4, 0, 0, false).
		AddItem(gui.friendList, 0, 4, 2, 1, 0, 0, false).
		AddItem(sendBtn, 2, 0, 1, 1, 0, 0, false).
		AddItem(exitBtn, 2, 2, 1, 1, 0, 0, false)

	gui.AddMsgInput()
	gui.LeaveChatHandler = func() {
		gui.app.Stop()
	}
}

// AddMsgInput adds the input field for typing in a chat message to the layout, this is needed
// because to clear an InputField in tview, we have to create a new InputField, so this code needs to run often
func (gui *ChatGUI) AddMsgInput() {
	gui.msgInput = tview.NewInputField()
	gui.msgInput.SetDoneFunc(gui.MsgInputHandler).
		SetBorder(true).
		SetTitle("Message").
		SetTitleAlign(tview.AlignLeft)

	gui.layout.AddItem(gui.msgInput, 1, 0, 1, 4, 0, 0, true)
	gui.app.SetFocus(gui.layout)
}

// MsgInputHandler is the key handler for the chat message input field
func (gui *ChatGUI) MsgInputHandler(key tcell.Key) {
	if key == tcell.KeyEnter {
		gui.SendDirectMessageHandler(gui.CurrentChatName, gui.msgInput.GetText())
		gui.layout.RemoveItem(gui.msgInput)
		gui.AddMsgInput()
	}
}

// Called when a friend is selected in the list
func (gui *ChatGUI) onFriendSelected(index int, name, secText string, scut rune) {
	gui.CurrentChatName = name

}

// KeyHandler is the keyboard input handler for the chat rooms interface
func (gui *ChatGUI) KeyHandler(key *tcell.EventKey) *tcell.EventKey {
	if key.Key() == tcell.KeyEsc {
		gui.LeaveChatHandler()
	}
	return key
}
