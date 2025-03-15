package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/app/chatview"
	"github.com/pjmessi/udp_chat/src/app/inputfield"
	"github.com/pjmessi/udp_chat/src/app/userlist"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

func getMainSectionComponent(chatView *tview.TextView, inputField *tview.InputField) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 3, 0, true)
}

func Run() {
	// Initialize state.
	appState := state.NewAppState()

	app := tview.NewApplication()

	messageSubmissionCh := make(chan inputfield.SubmissionEvent)
	userSelectionCh := make(chan userlist.SelectionEvent)

	// Initialize components.
	usersListSection := userlist.NewUserList(appState, userSelectionCh)
	chatViewSection := chatview.NewChatView()
	inputFieldSection := inputfield.NewInputField(appState, messageSubmissionCh)
	mainSection := getMainSectionComponent(chatViewSection.TextView, inputFieldSection.InputField)

	// Set up the overall layout.
	flex := tview.NewFlex().
		AddItem(usersListSection.List, 20, 0, true).
		AddItem(mainSection, 0, 1, false)

	// Handle new user selection event.
	go func(ch <-chan userlist.SelectionEvent) {
		for data := range ch {
			chatViewSection.UpdateView(appState)
			inputFieldSection.SetUserLabel(appState, &data.SelectedUser)
			app.SetFocus(inputFieldSection)
		}
	}(userSelectionCh)

	// Handle input field submission event.
	go func(ch <-chan inputfield.SubmissionEvent) {
		for data := range ch {
			if data.Status == inputfield.Sumibtted {
				chatViewSection.UpdateView(appState)
				app.SetFocus(inputFieldSection)
				continue
			}

			app.SetFocus(usersListSection)
		}
	}(messageSubmissionCh)

	// Set up key bindings for navigation
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Tab key to switch focus between user list and input field
		if event.Key() == tcell.KeyTab {
			if usersListSection.HasFocus() {
				app.SetFocus(inputFieldSection)
			} else {
				app.SetFocus(usersListSection)
			}
			return nil
		}

		// Global escape key handling (as a backup)
		if event.Key() == tcell.KeyEscape && inputFieldSection.HasFocus() {
			app.SetFocus(usersListSection)
			return nil
		}

		return event
	})

	// Run the application
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
