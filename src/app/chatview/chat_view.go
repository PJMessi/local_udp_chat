package chatview

import (
	"fmt"

	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type ChatView struct {
	*tview.TextView
}

func NewChatView() ChatView {
	chatView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
	chatView.SetBorder(true).SetTitle("Chat")

	return ChatView{
		TextView: chatView,
	}
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

		// Scroll to the end
		c.ScrollToEnd()
	}
}
