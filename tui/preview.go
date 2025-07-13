package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/water-sucks/optnix/option"
)

type PreviewModel struct {
	vp viewport.Model

	option       *option.NixosOption
	focused      bool
	lastRendered *option.NixosOption
}

func NewPreviewModel() PreviewModel {
	vp := viewport.New(0, 0)
	vp.SetHorizontalStep(1)

	return PreviewModel{
		vp: vp,
	}
}

func (m PreviewModel) SetHeight(height int) PreviewModel {
	m.vp.Height = height
	return m
}

func (m PreviewModel) SetWidth(width int) PreviewModel {
	m.vp.Width = width
	return m
}

func (m PreviewModel) SetFocused(focus bool) PreviewModel {
	m.focused = focus
	return m
}

func (m PreviewModel) SetOption(opt *option.NixosOption) PreviewModel {
	m.option = opt
	return m
}

var titleColor = color.New(color.Bold)

func (m PreviewModel) Update(msg tea.Msg) (PreviewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.option == nil {
				break
			}
			changeModeCmd := func() tea.Msg {
				return EvalValueStartMsg{Option: m.option.Name}
			}
			return m, changeModeCmd
		}
	case tea.WindowSizeMsg:
		// Force a re-render. The option string is cached otherwise,
		// and this can screw with the centered portion.
		m = m.ForceContentUpdate()
	}

	var cmd tea.Cmd
	if m.focused {
		m.vp, cmd = m.vp.Update(msg)
	}

	o := m.option

	// Do not re-render options if it has already been rendered before.
	// Setting content will reset the scroll counter, and rendering
	// an option is expensive.
	if o == m.lastRendered && o != nil {
		return m, cmd
	}

	m.vp.SetContent(m.renderOptionView())
	m.vp.GotoTop()

	m.lastRendered = o

	return m, cmd
}

func (m PreviewModel) ForceContentUpdate() PreviewModel {
	m.vp.SetContent(m.renderOptionView())
	m.vp.GotoTop()

	return m
}

func (m PreviewModel) renderOptionView() string {
	o := m.option

	sb := strings.Builder{}

	title := lipgloss.PlaceHorizontal(m.vp.Width, lipgloss.Center, titleColor.Sprint("Option Preview"))
	sb.WriteString(title)
	sb.WriteString("\n\n")

	if m.option == nil {
		sb.WriteString("\n  No option selected.")
		return sb.String()
	}

	sb.WriteString(o.PrettyPrint(nil))

	return sb.String()
}

func (m PreviewModel) View() string {
	if m.focused {
		m.vp.Style = focusedBorderStyle
	} else {
		m.vp.Style = inactiveBorderStyle
	}

	return m.vp.View()
}
