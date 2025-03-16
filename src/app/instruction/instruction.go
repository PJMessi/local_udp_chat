package instruction

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type InstructionView struct {
	*tview.TextView
}

func NewInstructionView(ctx context.Context, appState *state.AppState) InstructionView {
	tviewComp := tview.NewTextView().
		SetTextColor(tcell.ColorGray).
		SetText("Use arrow keys to move up/down, 'Enter' to select user/send a message, 'Esc' to get back to users list").
		SetTextAlign(tview.AlignCenter)

	instructionView := InstructionView{
		TextView: tviewComp,
	}

	return instructionView
}

