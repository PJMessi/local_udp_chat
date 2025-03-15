package userlist

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type UserList struct {
	*tview.List
}

func NewUserList(appState *state.AppState) UserList {
	component := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetMainTextColor(tcell.ColorWhite)

	userList := UserList{component}
	userList.RefreshList(appState)

	return userList
}

// This method accepts the function that is to be run when a list item is selected, and runs it on
// selection.
func (u *UserList) SetupHandler(handler func(index int)) {
	u.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		handler(index)
	})
}

// This method refreshes the list with potentially updated state data.
func (u *UserList) RefreshList(state *state.AppState) {
	u.Clear()

	selectedUser := state.SelectedUser

	for i, user := range state.GetUsers() {
		u.AddItem(user.Name, "", rune('1'+i%9), nil)
		if selectedUser != nil && user.Name == selectedUser.Name {
			u.SetCurrentItem(i)
		}
	}
}
