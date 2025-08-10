package cmdUtils

import (
	"os"
	"strings"

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
	varName := strings.ToUpper(prefix) + "_DEBUG_MODE"
	if os.Getenv(varName) == "" {
		return func() {}, nil
	}

	file, err := tea.LogToFile(prefix+".debug.log", prefix)

	return func() {
		if err != nil || file == nil {
			return
		}
		_ = file.Close()
	}, err
}
