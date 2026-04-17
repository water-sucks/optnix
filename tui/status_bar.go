package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBarModel struct {
	defaultText string
	text        string
	kind        NotificationKind
	id          int
	width       int
}

func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		defaultText: "For basic help, press Ctrl-G.",
	}
}

func (m StatusBarModel) SetWidth(width int) StatusBarModel {
	m.width = width
	return m
}

func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case NotificationMsg:
		m.text = msg.Message
		m.kind = msg.Kind
		m.id++
		id := m.id
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return ClearNotificationMsg{ID: id}
		})
	case ClearNotificationMsg:
		if msg.ID == m.id {
			m.text = ""
		}
		return m, nil
	}
	return m, nil
}

func (m StatusBarModel) View() string {
	text, style := m.defaultText, hintStyle
	if m.text != "" {
		text = m.text
		if m.kind == NotificationError {
			style = errorHintStyle
		}
	}
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, style.Render(text))
}
