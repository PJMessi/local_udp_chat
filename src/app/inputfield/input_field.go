package inputfield

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type MessageStatus int

const (
	Sumibtted           MessageStatus = iota
	SubmissionCancelled               = iota
)

type SubmissionEvent struct {
	Status MessageStatus
}

type InputField struct {
	*tview.InputField
}

func NewInputField(appState *state.AppState, ch chan SubmissionEvent) InputField {
	tviewInputField := tview.NewInputField().SetLabel("> ").SetFieldWidth(0)
	tviewInputField.SetBorder(true).SetTitle("Message")
	inputField := InputField{InputField: tviewInputField}
	inputField.setupHandler(appState, ch)
	return inputField
}

func (i *InputField) SetUserLabel(_ *state.AppState, selectedUser *state.User) {
	i.SetLabel(fmt.Sprintf("[To: %s] > ", selectedUser.Name))
}

func (i *InputField) setupHandler(appState *state.AppState, ch chan<- SubmissionEvent) {
	i.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		text := i.GetText()
		if strings.TrimSpace(text) != "" {
			selectedUser := appState.GetSelectedUser()
			selectedUser.AppendMessage(appState, text, state.MessageSourceSelf)
			i.SetText("")
			ch <- SubmissionEvent{}
		}
	})
}

// 			// Simulate a reply after a brief delay
// 			go func() {
// 				time.Sleep(time.Second)
// 				app.QueueUpdateDraw(func() {
// 					selectedUser.AppendMessage(appState, fmt.Sprintf("You said: '%s'", text), state.MessageSourceUser)
// 					// Update chat view.
// 					chatViewSection.UpdateView(appState)
// 				})
// 			}()
