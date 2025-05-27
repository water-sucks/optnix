package cmdUtils

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type ErrorWithHint struct {
	Msg  string
	Hint string
}

func (e ErrorWithHint) Error() string {
	msg := e.Msg
	if e.Hint != "" {
		msg += "\n\n" + color.YellowString("hint: %v", e.Hint)
	}
	msg += "\n\nTry 'optnix --help' for more information."

	return msg
}

func ConfigureBubbleTeaLogger(prefix string) (func(), error) {
	if os.Getenv("OPTNIX_DEBUG_MODE") == "" {
		return func() {}, nil
	}

	file, err := tea.LogToFile("debug.log", prefix)

	return func() {
		if err != nil || file == nil {
			return
		}
		_ = file.Close()
	}, err
}
