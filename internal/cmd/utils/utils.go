package cmdUtils

import "github.com/fatih/color"

type ArgParseError struct {
	Msg  string
	Hint string
}

func (e ArgParseError) Error() string {
	msg := e.Msg
	if e.Hint != "" {
		msg += "\n\n" + color.YellowString("hint: %v", e.Hint)
	}
	msg += "\n\nTry 'optnix --help' for more information."

	return msg
}
