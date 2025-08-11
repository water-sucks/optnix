package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/water-sucks/optnix/option"
)

var (
	ansiRed     = lipgloss.ANSIColor(termenv.ANSIRed)
	ansiYellow  = lipgloss.ANSIColor(termenv.ANSIYellow)
	ansiGreen   = lipgloss.ANSIColor(termenv.ANSIGreen)
	ansiWhite   = lipgloss.ANSIColor(termenv.ANSIBrightWhite)
	ansiBlue    = lipgloss.ANSIColor(termenv.ANSIBlue)
	ansiCyan    = lipgloss.ANSIColor(termenv.ANSICyan)
	ansiMagenta = lipgloss.ANSIColor(termenv.ANSIMagenta)

	itemStyle         = lipgloss.NewStyle().MarginLeft(4).PaddingLeft(1).Border(lipgloss.NormalBorder(), false, false, false, true)
	currentItemStyle  = lipgloss.NewStyle().MarginLeft(4).PaddingLeft(1).Foreground(ansiGreen).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(ansiGreen)
	selectedItemStyle = lipgloss.NewStyle().MarginLeft(4).PaddingLeft(1).Foreground(ansiYellow).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(ansiYellow)
	errorTextStyle    = lipgloss.NewStyle().Foreground(ansiRed)
	attrStyle         = lipgloss.NewStyle().Foreground(ansiCyan)
	boldStyle         = lipgloss.NewStyle().Bold(true)
	italicStyle       = lipgloss.NewStyle().Italic(true)
)

type LoadScopeStartMsg option.Scope

type LoadScopeFinishedMsg struct {
	Name    string
	Options option.NixosOptionSource
	Err     error
}

type ChangeScopeMsg struct {
	Name    string
	Options option.NixosOptionSource
}

type scopeItem struct {
	Scope    option.Scope
	Selected bool
}

func (i scopeItem) FilterValue() string {
	scope := i.Scope
	return fmt.Sprintf("%v %v", scope.Name, scope.Description)
}

type scopeItemDelegate struct{}

func (d scopeItemDelegate) Height() int { return 2 }

func (d scopeItemDelegate) Spacing() int { return 1 }

func (d scopeItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

var selectedText = italicStyle.Render(" (selected)")

func (d scopeItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(scopeItem)
	if !ok {
		return
	}

	s := i.Scope

	var str string
	if !i.Selected {
		str = boldStyle.Render(s.Name)
	} else {
		str = boldStyle.Render(fmt.Sprintf("%v%v", s.Name, selectedText))
	}

	str += fmt.Sprintf("\n%s :: %s", attrStyle.Render("Description"), s.Description)

	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return currentItemStyle.Render(strings.Join(s, " "))
		}
	} else if i.Selected {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

type SelectScopeModel struct {
	vp      viewport.Model
	spinner spinner.Model
	list    list.Model

	width  int
	height int

	scopes        []option.Scope
	selectedScope string

	loading bool
	err     error
}

func NewSelectScopeModel(scopes []option.Scope, selectedScope string) SelectScopeModel {
	items := make([]list.Item, len(scopes))
	for i, s := range scopes {
		selected := s.Name == selectedScope

		items[i] = scopeItem{
			Scope:    s,
			Selected: selected,
		}
	}

	l := list.New(items, scopeItemDelegate{}, 0, 0)

	l.Title = "Available Scopes"

	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2).Background(ansiRed).Foreground(ansiWhite)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	l.Styles.StatusBar = lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1).Foreground(ansiMagenta)

	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(ansiBlue).Bold(true).PaddingLeft(2)
	l.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(ansiBlue)
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(ansiBlue)
	l.Styles.StatusBarActiveFilter = lipgloss.NewStyle().Foreground(ansiBlue)
	l.Styles.StatusBarFilterCount = lipgloss.NewStyle().Foreground(ansiBlue)

	vp := viewport.New(0, 0)
	vp.SetHorizontalStep(1)
	vp.Style = focusedBorderStyle

	sp := spinner.New()
	sp.Spinner = spinner.Line
	sp.Style = spinnerStyle

	return SelectScopeModel{
		list:    l,
		vp:      vp,
		spinner: sp,

		scopes:        scopes,
		selectedScope: selectedScope,
	}
}

func (m SelectScopeModel) Update(msg tea.Msg) (SelectScopeModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if m.err != nil {
			m.err = nil
			break
		}

		switch msg.String() {
		case "q", "esc":
			return m, func() tea.Msg {
				return ChangeViewModeMsg(ViewModeSearch)
			}

		case "enter":
			item := m.list.SelectedItem().(scopeItem)
			if item.Selected {
				break
			}

			return m, func() tea.Msg {
				return LoadScopeStartMsg(item.Scope)
			}
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)

		m.width = msg.Width - 4
		m.height = msg.Height - 4

		m.vp.Width = m.width
		m.vp.Height = m.height

		return m, nil

	case LoadScopeStartMsg:
		m.loading = true
		m.err = nil

		cmds = append(cmds, m.spinner.Tick)
		cmds = append(cmds, func() tea.Msg {
			loaded, err := msg.Loader()
			return LoadScopeFinishedMsg{
				Name:    msg.Name,
				Options: loaded,
				Err:     err,
			}
		})

		return m, tea.Batch(cmds...)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case LoadScopeFinishedMsg:
		m.loading = false
		m.err = msg.Err

		if m.err == nil {
			return m, func() tea.Msg {
				return ChangeScopeMsg{
					Name:    msg.Name,
					Options: msg.Options,
				}
			}
		}
	}

	if m.loading {
		m.vp.SetContent("Loading..." + m.spinner.View())
	} else if m.err != nil {
		m.vp.SetContent(errorTextStyle.Render(fmt.Sprintf("error loading scope: %v\n\nPress any key to go back.\n", m.err)))
	}

	var vpCmd tea.Cmd
	m.vp, vpCmd = m.vp.Update(msg)
	cmds = append(cmds, vpCmd)

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m SelectScopeModel) View() string {
	if m.loading || m.err != nil {
		return m.vp.View()
	}

	return "\n\n" + m.list.View()
}
