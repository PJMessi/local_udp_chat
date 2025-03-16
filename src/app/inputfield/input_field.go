package inputfield

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/network/msgsender"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/pjmessi/udp_chat/src/utils/ctxutils"
	"github.com/rivo/tview"
)

type MessageStatus int

const (
	Submitted           MessageStatus = iota
	SubmissionCancelled               = iota
)

type SubmissionEvent struct {
	Status MessageStatus
}

type InputField struct {
	*tview.InputField
}

func NewInputField(ctx context.Context, appState *state.AppState, ch chan SubmissionEvent) InputField {
	tviewInputField := tview.
		NewInputField().
		SetFieldBackgroundColor(tcell.ColorNone).
		SetLabel("> ").
		SetFieldWidth(0)
	tviewInputField.SetBorder(true).SetTitle("Message")
	inputField := InputField{InputField: tviewInputField}
	inputField.setupHandler(ctx, appState, ch)
	return inputField
}

func (i *InputField) SetUserLabel(_ *state.AppState, selectedUser *state.User) {
	i.SetLabel(fmt.Sprintf("[To: %s] > ", selectedUser.Name))
}

func (i *InputField) setupHandler(ctx context.Context, appState *state.AppState, ch chan<- SubmissionEvent) {
	logger := ctxutils.ExtractLogger(ctx)

	i.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		text := i.GetText()
		if strings.TrimSpace(text) != "" {
			selectedUser := appState.GetSelectedUser()
			selectedUser.AppendMessage(appState, text, state.MessageSourceSelf)
			i.SetText("")

			err := msgsender.SendMessage(ctx, appState.SelfUsername, selectedUser, text)
			if err != nil {
				logger.Error("error while sending message", "error", err)
				return
			}

			ch <- SubmissionEvent{}
		}
	})
}
