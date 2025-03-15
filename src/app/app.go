package app

import (
	"fmt"
	"strings"
	"time"

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

	// Initialize components.
	usersListSection := userlist.NewUserList(appState)
	chatViewSection := chatview.NewChatView()
	inputFieldSection := inputfield.NewInputField()
	mainSection := getMainSectionComponent(chatViewSection.TextView, inputFieldSection.InputField)

	// usersListSection.List
	// selectedUser := 0

	// Set up the overall layout.
	flex := tview.NewFlex().
		AddItem(usersListSection.List, 20, 0, true).
		AddItem(mainSection, 0, 1, false)

	// Marks the user with the given index as selected, update the chat view with the user's message,
	// and updates the label in the input section, and focus the input section.
	updateChatView := func(index int) {
		// Update selected user.
		user := appState.SelectUser(uint(index))
		// Update chat view.
		chatViewSection.UpdateView(appState)
		// Update input section.
		inputFieldSection.SetUserLabel(appState, user)
		// Set focus to input section.
		app.SetFocus(inputFieldSection)
	}

	usersListSection.SetupHandler(updateChatView)

	// Handle input field submit
	inputFieldSection.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := inputFieldSection.GetText()
			if strings.TrimSpace(text) != "" {
				selectedUser := appState.GetSelectedUser()

				// Add the message to the current user's chat
				selectedUser.AppendMessage(appState, text, state.MessageSourceSelf)

				// Simulate a reply after a brief delay
				go func() {
					time.Sleep(time.Second)
					app.QueueUpdateDraw(func() {
						selectedUser.AppendMessage(appState, fmt.Sprintf("You said: '%s'", text), state.MessageSourceUser)
						// Update chat view.
						chatViewSection.UpdateView(appState)
					})
				}()

				// Clear input and update view
				inputFieldSection.SetText("")
				// Update chat view.
				chatViewSection.UpdateView(appState)
				// Set focus to input section.
				app.SetFocus(inputFieldSection)
			}
		} else if key == tcell.KeyEscape {
			// When Escape is pressed, go back to user selection
			app.SetFocus(usersListSection)
		}
	})

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

	// Initial view setup - select the first user
	updateChatView(0)

	// Run the application
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
