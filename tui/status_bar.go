package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBarModel struct {
	defaultText string
	flashText   string
	flashActive bool
	flashID     int
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
	case CopiedToClipboardMsg:
		m.flashActive = true
		m.flashText = "Copied to clipboard!"
		m.flashID++
		id := m.flashID
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return ClearClipboardFlashMsg{id: id}
		})
	case ClearClipboardFlashMsg:
		if msg.id == m.flashID {
			m.flashActive = false
			m.flashText = ""
		}
		return m, nil
	}
	return m, nil
}

func (m StatusBarModel) View() string {
	text := m.defaultText
	if m.flashActive {
		text = m.flashText
	}
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, hintStyle.Render(text))
}
