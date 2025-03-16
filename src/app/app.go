package app

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/app/chatview"
	"github.com/pjmessi/udp_chat/src/app/inputfield"
	"github.com/pjmessi/udp_chat/src/app/userlist"
	"github.com/pjmessi/udp_chat/src/logger"
	"github.com/pjmessi/udp_chat/src/network/broadcaster"
	"github.com/pjmessi/udp_chat/src/network/detector"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/pjmessi/udp_chat/src/utils/ctxutils"
	"github.com/rivo/tview"
)

func getMainSectionComponent(chatView *tview.TextView, inputField *tview.InputField) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 3, 0, true)
}

func Run(username string) {
	logger := logger.NewLogger()
	ctx := ctxutils.AttachLogger(context.Background(), logger)
	appState := state.NewAppState()
	appState.SetUsername(username)

	// Broadcast self continiuously.
	go func() {
		err := broadcaster.Broadcast(ctx, appState.SelfUsername)
		if err != nil {
			logger.Error("error while broadcasting self", "error", err)
		}
	}()

	app := tview.NewApplication()
	messageSubmissionCh := make(chan inputfield.SubmissionEvent)
	userSelectionCh := make(chan userlist.SelectionEvent)
	searchCompleteCh := make(chan detector.SearchCompleteEvent)
	newMessageCh := make(chan detector.NewMessageEvent)

	// Listen for incoming requests continuously
	go func() {
		err := detector.Listen(ctx, appState, searchCompleteCh, newMessageCh)
		if err != nil {
			panic(err)
		}
	}()

	// Initialize components.
	usersListSection := userlist.NewUserList(appState, userSelectionCh)
	chatViewSection := chatview.NewChatView()
	inputFieldSection := inputfield.NewInputField(ctx, appState, messageSubmissionCh)
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

	// Handle new message event.
	go func(ch <-chan detector.NewMessageEvent) {
		for data := range ch {
			app.QueueUpdateDraw(func() {
				chatViewSection.UpdateView(appState)
				// Only focus the input field, if the message is from the currently selected user.
				if appState.SelectedUser != nil && appState.SelectedUser.Name == data.User.Name {
					app.SetFocus(inputFieldSection)
				}
			})
		}
	}(newMessageCh)

	// Handle input field submission event.
	go func(ch <-chan inputfield.SubmissionEvent) {
		for data := range ch {
			app.QueueUpdateDraw(func() {
				if data.Status == inputfield.SubmissionCancelled {
					app.SetFocus(usersListSection)
					return
				}

				chatViewSection.UpdateView(appState)
				app.SetFocus(inputFieldSection)
			})
		}
	}(messageSubmissionCh)

	// Handle search results.
	go func(ch <-chan detector.SearchCompleteEvent) {
		for range ch {
			app.QueueUpdateDraw(func() {
				usersListSection.RefreshList(appState)
			})
		}
	}(searchCompleteCh)

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
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
