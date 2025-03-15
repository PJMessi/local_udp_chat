package userlist

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type UserSelectionEvent struct {
	SelectedUser state.User
}

type UserList struct {
	*tview.List
}

func NewUserList(appState *state.AppState, ch chan UserSelectionEvent) UserList {
	component := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetMainTextColor(tcell.ColorWhite)

	userList := UserList{component}
	userList.RefreshList(appState)
	userList.setupHandler(appState, ch)

	return userList
}

func (u *UserList) setupHandler(appState *state.AppState, ch chan<- UserSelectionEvent) {
	u.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		user := appState.SelectUser(uint(index))
		if user != nil {
			ch <- UserSelectionEvent{
				SelectedUser: *user,
			}
		}
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
