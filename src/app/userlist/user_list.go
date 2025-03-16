package userlist

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type SelectionEvent struct {
	SelectedUser state.User
}

type UserList struct {
	*tview.List
}

func NewUserList(ctx context.Context, appState *state.AppState, ch chan SelectionEvent) UserList {
	component := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetMainTextColor(tcell.ColorWhite)

	userList := UserList{component}
	userList.RefreshList(ctx, appState)
	userList.setupHandler(ctx, appState, ch)

	return userList
}

func (u *UserList) setupHandler(_ context.Context, appState *state.AppState, ch chan<- SelectionEvent) {
	u.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		user := appState.SelectUser(uint(index))
		if user != nil {
			ch <- SelectionEvent{
				SelectedUser: *user,
			}
		}
	})
}

func (u *UserList) RefreshList(ctx context.Context, state *state.AppState) {
	u.Clear()

	selectedUser := state.SelectedUser
	users := state.GetUsers()

	if len(users) > 0 {
		for i, user := range users {
			u.AddItem(user.Name, "", rune('1'+i%9), nil)
			if selectedUser != nil && user.Name == selectedUser.Name {
				u.SetCurrentItem(i)
			}
		}
	}
}
