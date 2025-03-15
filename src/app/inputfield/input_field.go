package inputfield

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/pjmessi/udp_chat/src/state"
	"github.com/rivo/tview"
)

type InputField struct {
	*tview.InputField
}

func NewInputField() InputField {
	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)
	inputField.SetBorder(true).SetTitle("Message")

	return InputField{
		InputField: inputField,
	}
}

func (i *InputField) SetUserLabel(_ *state.AppState, selectedUser *state.User) {
	i.SetLabel(fmt.Sprintf("[To: %s] > ", selectedUser.Name))
}

// Sets up the input field to execute the provided handler function when done.
func (i *InputField) SetupHandler(handler func(tcell.Key)) {
	i.SetDoneFunc(handler)
}
