package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/grigagod/chat-example/sdb"
	"github.com/rivo/tview"
)

// ChatGUI contains the widgets/state for the chat main window view
type ChatGUI struct {
	*GUI

	sendDirectMessageHandler    func(friendname string, msg string)
	acceptFriendRequestHandler  func(friendname string)
	declineFriendRequestHandler func(friendname string)
	chatInfoHandler             func()
	leaveChatHandler            func()
	selectedFriendName          string
	getMsgList                  func(friendname string) []*sdb.Message
	// friends					 map[string]*big.Int // tmp
	invitesListUpdater *time.Ticker

	layout           *tview.Grid
	friendsListView  *tview.List
	requestsListView *tview.List
	msgView          *tview.TextView
	msgInput         *tview.InputField
	addFriendBtn     *tview.Button
	checkInvitesBtn  *tview.Button

	focusableElements []tview.Primitive
	focusedIndex      int
}

// Create initializes the widgets in the chat GUI
func (gui *ChatGUI) Create() {
	gui.friendsListView = tview.NewList()
	gui.friendsListView.
		SetSelectedFunc(gui.onFriendSelected).
		SetTitle("Friends").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)

	gui.msgView = tview.NewTextView()
	gui.msgView.SetDynamicColors(true).SetScrollable(true).
		SetBorder(true).
		SetTitle("Chat")

	gui.requestsListView = tview.NewList()
	gui.requestsListView.
		SetSelectedFunc(gui.onRequestSelected).
		SetTitle("Requests").
		SetBorder(true).
		SetTitleAlign(tview.AlignCenter)

	gui.addFriendBtn = tview.NewButton("Add friend")

	gui.layout = tview.NewGrid()
	gui.layout.SetRows(0, 6, 0, 3, 1).
		SetColumns(20, 1, 20, 0, 30).
		AddItem(gui.msgView, 0, 0, 3, 4, 0, 0, false).
		AddItem(gui.friendsListView, 0, 4, 2, 1, 0, 0, false).
		AddItem(gui.requestsListView, 2, 4, 2, 1, 0, 0, false).
		AddItem(gui.addFriendBtn, 4, 4, 1, 1, 0, 0, false)

	gui.leaveChatHandler = func() {
		gui.app.Stop()
	}

	gui.focusableElements = []tview.Primitive{
		gui.msgView,
		gui.friendsListView,
		gui.requestsListView,
		gui.addFriendBtn,
	}

	gui.AddMsgInput()
	gui.focusedIndex = 1

}

// AddMsgInput adds the input field for typing in a chat message to the layout, this is needed
// because to clear an InputField in tview, we have to create a new InputField, so this code needs to run often
func (gui *ChatGUI) AddMsgInput() {
	gui.msgInput = tview.NewInputField()
	gui.msgInput.SetDoneFunc(gui.MsgInputHandler).
		SetBorder(true).
		SetTitle("Message").
		SetTitleAlign(tview.AlignLeft)

	gui.focusableElements = append(gui.focusableElements, gui.msgInput)

	gui.layout.AddItem(gui.msgInput, 3, 0, 1, 4, 0, 0, true)
	//	gui.app.SetFocus(gui.layout)
}

// MsgInputHandler is the key handler for the chat message input field
func (gui *ChatGUI) MsgInputHandler(key tcell.Key) {
	if key == tcell.KeyEnter && gui.selectedFriendName != "" {
		for _, primitive := range gui.focusableElements {
			if primitive == gui.msgInput {
				gui.sendDirectMessageHandler(gui.selectedFriendName, gui.msgInput.GetText())

				gui.msgInput.SetText("")
			}
		}
	}
}

func formatChatMessage(sender, message string, timestamp int64) []byte {
	var buf bytes.Buffer
	log.Println("start formatting")
	tm := time.Unix(timestamp/1000, 0)
	buf.WriteString(fmt.Sprintf("[dimgray]%02d-%02d %02d:%02d[white]", tm.Day(), tm.Month(), tm.Hour(), tm.Minute()))
	buf.WriteString(" [blue]<")
	buf.WriteString(string(sender))
	buf.WriteString("> [white]")
	buf.WriteString(message)
	buf.WriteRune('\n')

	return buf.Bytes()
}

func (gui *ChatGUI) DisplayMessage(username, msg string, timestamp int64) {
	gui.app.QueueUpdateDraw(func() {
		fmtMsg := formatChatMessage(username, msg, timestamp)
		gui.msgView.Write(fmtMsg)
		gui.msgView.ScrollToEnd()
	})
}

func (gui *ChatGUI) InstantDisplayMessage(username, msg string, timestamp int64) {
	fmtMsg := formatChatMessage(username, msg, timestamp)
	gui.msgView.Write(fmtMsg)
	gui.msgView.ScrollToEnd()
}

// Called when a friend is selected in the list
func (gui *ChatGUI) onFriendSelected(index int, name, secText string, scut rune) {
	if gui.selectedFriendName != name {
		gui.selectedFriendName = name
		gui.loadMsgHistory()
	}
}

func (gui *ChatGUI) onRequestSelected(index int, name, secText string, scut rune) {
	gui.acceptFriendRequestHandler(name)
	gui.removeCurrentRequest()
	log.Println(name)
}

// KeyHandler is the keyboard input handler for the chat rooms interface
func (gui *ChatGUI) KeyHandler(key *tcell.EventKey) *tcell.EventKey {
	if key.Key() == tcell.KeyEsc {
		gui.app.Stop()
	}
	if key.Key() == tcell.KeyTab {
		gui.focusedIndex++
		if gui.focusedIndex == len(gui.focusableElements) {
			gui.focusedIndex = 0
		}
		gui.app.SetFocus(gui.focusableElements[gui.focusedIndex])
		return nil
	} else if key.Key() == tcell.KeyEnter {
		switch gui.app.GetFocus() {
		case gui.addFriendBtn:
			gui.chatInfoHandler()
			return nil
		}
	}
	return key
}

func (gui *ChatGUI) addToFriendList(friend string) {
	gui.friendsListView.AddItem(friend, "", 0, nil)
}

func (gui *ChatGUI) addToRequestsList(request string) {
	gui.requestsListView.AddItem(request, "", 0, nil)
}

func (gui *ChatGUI) removeCurrentRequest() {
	gui.requestsListView.RemoveItem(gui.requestsListView.GetCurrentItem())
}

func (gui *ChatGUI) loadMsgHistory() {
	gui.msgView.Clear()

	messages := gui.getMsgList(gui.selectedFriendName)
	for _, msg := range messages {
		gui.InstantDisplayMessage(msg.SenderName, msg.Message, msg.Timestamp)
	}
	gui.msgView.ScrollToEnd()
}
