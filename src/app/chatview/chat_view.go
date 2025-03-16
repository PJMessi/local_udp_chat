package chatview

import (
	"context"
	"fmt"

	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type ChatView struct {
	*tview.TextView
}

func NewChatView(ctx context.Context, appState *state.AppState) ChatView {
	tviewComp := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
	tviewComp.SetBorder(true).SetTitle("Chat")

	chatView := ChatView{
		TextView: tviewComp,
	}

	chatView.UpdateView(appState)

	return chatView
}

func (c *ChatView) UpdateView(appState *state.AppState) {
	// Clear and update chat view.
	c.Clear()

	selectedUser := appState.GetSelectedUser()
	if selectedUser != nil {
		for _, msg := range selectedUser.Messages {
			timeStr := msg.Timestamp.Format("15:04")

			if msg.Source == state.MessageSourceSelf {
				fmt.Fprintf(c, "[green]%s [%s]:[white] %s\n", "You", timeStr, msg.Content)

			} else {
				fmt.Fprintf(c, "[blue]%s [%s]:[white] %s\n", selectedUser.Name, timeStr, msg.Content)
			}
		}

	} else {
		fmt.Fprintf(c, "select a user")
	}

	// Scroll to the end
	c.ScrollToEnd()
}
