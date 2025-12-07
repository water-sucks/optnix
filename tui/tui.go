package tui

import (
	"fmt"
	"regexp"
	"slices"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/sahilm/fuzzy"
	cmdUtils "github.com/water-sucks/optnix/internal/cmd/utils"
	"github.com/water-sucks/optnix/internal/utils"
	"github.com/water-sucks/optnix/option"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)

	inactiveBorderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	focusedBorderStyle  = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.ANSIColor(termenv.ANSIMagenta))
	titleRuleStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.ANSIColor(termenv.ANSIWhite)).
			BorderTop(true).
			BorderRight(false).
			BorderBottom(false).
			BorderLeft(false)

	marginStyle = lipgloss.NewStyle().Margin(2, 2, 0, 2)
	hintStyle   = lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(termenv.ANSIYellow)) // Soft gray

)

type Model struct {
	focus FocusArea
	mode  ViewMode

	options              option.NixosOptionSource
	enableScopeSwitching bool

	filtered []fuzzy.Match
	minScore int64

	width  int
	height int

	search      SearchBarModel
	results     ResultListModel
	preview     PreviewModel
	selectScope SelectScopeModel
	eval        EvalValueModel
	help        HelpModel
}

type ViewMode int

const (
	ViewModeSearch = iota
	ViewModeSelectScope
	ViewModeEvalValue
	ViewModeHelp
)

type ChangeViewModeMsg ViewMode

type FocusArea int

const (
	FocusAreaResults FocusArea = iota
	FocusAreaPreview
)

func NewModel(
	scopes []option.Scope,
	selectedScope string,
	minScore int64,
	debounceTime int64,
	initialInput string,
) (*Model, error) {
	var scope *option.Scope
	for _, s := range scopes {
		if selectedScope == s.Name {
			scope = &s
			break
		}
	}

	if scope == nil {
		return nil, fmt.Errorf("scope '%v' not found in configuration", selectedScope)
	}

	options, err := scope.Loader()
	if err != nil {
		return nil, err
	}

	preview := NewPreviewModel()
	search := NewSearchBarModel(len(options), debounceTime).
		SetFocused(true).
		SetValue(initialInput)
	results := NewResultListModel(options, scope.Name).
		SetFocused(true)
	selectScope := NewSelectScopeModel(scopes, scope.Name)
	eval := NewEvalValueModel(scope.Evaluator)
	help := NewHelpModel()

	return &Model{
		mode:  ViewModeSearch,
		focus: FocusAreaResults,

		options:              options,
		enableScopeSwitching: len(scopes) > 1,

		minScore: minScore,

		results:     results,
		preview:     preview,
		search:      search,
		selectScope: selectScope,
		eval:        eval,
		help:        help,
	}, nil
}

func (m Model) Init() tea.Cmd {
	if m.search.Value() != "" {
		return func() tea.Msg {
			return RunSearchMsg{
				Query: m.search.Value(),
				Mode:  SearchModeFuzzy,
			}
		}
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.mode != ViewModeEvalValue {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m = m.updateWindowSize(msg.Width, msg.Height)

		// Always forward resize events to components that need them.
		m.eval, _ = m.eval.Update(msg)
		m.help, _ = m.help.Update(msg)
		m.selectScope, _ = m.selectScope.Update(msg)

		return m, nil

	case ChangeViewModeMsg:
		m.mode = ViewMode(msg)

	case EvalValueStartMsg:
		m.mode = ViewModeEvalValue

	case ChangeScopeMsg:
		m.mode = ViewModeSearch
		m.options = msg.Options
	}

	switch m.mode {
	case ViewModeSearch:
		return m.updateSearch(msg)
	case ViewModeEvalValue:
		var evalCmd tea.Cmd
		m.eval, evalCmd = m.eval.Update(msg)
		return m, evalCmd
	case ViewModeSelectScope:
		var selectModeCmd tea.Cmd
		m.selectScope, selectModeCmd = m.selectScope.Update(msg)
		return m, selectModeCmd
	case ViewModeHelp:
		var helpCmd tea.Cmd
		m.help, helpCmd = m.help.Update(msg)
		return m, helpCmd
	}

	return m, nil
}

func (m Model) updateSearch(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m = m.toggleFocus()

		case "ctrl+g":
			return m, func() tea.Msg {
				return ChangeViewModeMsg(ViewModeHelp)
			}
		case "ctrl+o":
			if !m.enableScopeSwitching {
				return m, nil
			}

			return m, func() tea.Msg {
				return ChangeViewModeMsg(ViewModeSelectScope)
			}
		}
	case RunSearchMsg:
		m = m.runSearch(msg.Query, msg.Mode)
		m.search = m.search.SetResultCount(len(m.filtered))
	}

	var cmds []tea.Cmd

	var searchCmd tea.Cmd
	m.search, searchCmd = m.search.Update(msg)
	cmds = append(cmds, searchCmd)

	var resultsCmd tea.Cmd
	m.results, resultsCmd = m.results.Update(msg)
	cmds = append(cmds, resultsCmd)

	selectedOption := m.results.GetSelectedOption()
	m.preview = m.preview.SetOption(selectedOption)

	var previewCmd tea.Cmd
	m.preview, previewCmd = m.preview.Update(msg)
	cmds = append(cmds, previewCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) runSearch(query string, mode SearchMode) Model {
	m.results = m.results.SetSearchError(nil)

	if len(query) == 0 {
		m.filtered = nil
		m.results = m.results.
			SetQuery(query).
			SetResultList(m.filtered).
			SetSelectedIndex(len(m.filtered) - 1)
		return m
	}

	var matches []fuzzy.Match
	switch mode {
	case SearchModeFuzzy:
		allMatches := fuzzy.FindFrom(query, m.options)
		matches = utils.FilterMinimumScoreMatches(allMatches, m.minScore)

		// Reverse the filtered match list, since we want more relevant
		// options at the bottom of the screen.
		slices.Reverse(matches)
	case SearchModeRegex:
		expr, err := regexp.Compile(query)
		if err != nil {
			m.results = m.results.SetSearchError(err)
			return m
		}
		matches = regexSearch(m.options, expr)
	default:
		panic("unhandled search mode")
	}

	m.filtered = matches

	m.results = m.results.
		SetQuery(query).
		SetResultList(m.filtered).
		SetSelectedIndex(len(m.filtered) - 1)

	return m
}

func regexSearch(options option.NixosOptionSource, expr *regexp.Regexp) []fuzzy.Match {
	var matches []fuzzy.Match

	for i, o := range options {
		matchedCaptureRanges := expr.FindAllStringSubmatchIndex(o.Name, -1)
		if len(matchedCaptureRanges) == 0 {
			continue
		}

		m := fuzzy.Match{}

		m.Index = i

		for _, capture := range matchedCaptureRanges {
			start, end := capture[0], capture[1]-1

			capturedIndices := make([]int, end-start+1)
			for j := range capturedIndices {
				capturedIndices[j] = start + j
			}

			m.Str = o.Name
			m.MatchedIndexes = append(m.MatchedIndexes, capturedIndices...)
			m.Score = calculateRegexScore(o.Name, capturedIndices)
		}

		matches = append(matches, m)
	}

	slices.SortFunc(matches, func(a, b fuzzy.Match) int {
		return a.Score - b.Score
	})

	return matches
}

func calculateRegexScore(str string, matchedIndexes []int) int {
	if len(matchedIndexes) == 0 {
		return 0
	}

	const (
		consecutiveBonus  = 5
		wordBoundaryBonus = 10
	)

	// Base score: 1 per matched character
	score := len(matchedIndexes)

	// Bonus for consecutive matches
	for i := 1; i < len(matchedIndexes); i++ {
		if matchedIndexes[i] == matchedIndexes[i-1]+1 {
			score += consecutiveBonus
		}
	}

	// Bonus for matches at start of word boundaries
	for _, idx := range matchedIndexes {
		if idx == 0 {
			// First character is automatically a word boundary
			score += wordBoundaryBonus
		} else {
			// Matched characters after non-alphanumeric ones
			// also are word boundaries
			char := rune(str[idx-1])
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				score += wordBoundaryBonus
			}
		}
	}

	// Normalize by string length
	return score * 100 / len(str)
}

type RunSearchMsg struct {
	Query string
	Mode  SearchMode
}

func (m Model) toggleFocus() Model {
	switch m.focus {
	case FocusAreaResults:
		m.focus = FocusAreaPreview

		m.results = m.results.SetFocused(false)
		m.search = m.search.SetFocused(false)
		m.preview = m.preview.SetFocused(true)
	case FocusAreaPreview:
		m.focus = FocusAreaResults

		m.results = m.results.SetFocused(true)
		m.search = m.search.SetFocused(true)
		m.preview = m.preview.SetFocused(false)
	}

	return m
}

func (m Model) updateWindowSize(width, height int) Model {
	m.width = width
	m.height = height

	usableWidth := width - 4   // 2 left + 2 right margins
	usableHeight := height - 2 // 2 top margin

	searchHeight := 3

	halfWidth := usableWidth / 2

	m.results = m.results.
		SetWidth(halfWidth - 2). // 1 border each side
		SetHeight(usableHeight - searchHeight - 2)

	m.search = m.search.
		SetWidth(halfWidth - 2).
		SetHeight(searchHeight)

	m.preview = m.preview.
		SetWidth(halfWidth - 2).
		SetHeight(usableHeight - 2)

	return m
}

func (m Model) View() string {
	switch m.mode {
	case ViewModeSelectScope:
		return marginStyle.Render(m.selectScope.View())
	case ViewModeEvalValue:
		return marginStyle.Render(m.eval.View())
	case ViewModeHelp:
		return marginStyle.Render(m.help.View())
	}

	results := m.results.View()
	search := m.search.View()
	preview := m.preview.View()

	left := lipgloss.JoinVertical(lipgloss.Top, results, search)
	main := lipgloss.JoinHorizontal(lipgloss.Top, left, preview)

	hint := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, hintStyle.Render("For basic help, press Ctrl-G."))

	return lipgloss.JoinVertical(
		lipgloss.Top,
		marginStyle.Render(main),
		hint,
	)
}

type OptionTUIArgs struct {
	Scopes            []option.Scope
	SelectedScopeName string
	MinScore          int64
	DebounceTime      int64
	InitialInput      string
	LogFileName       string
}

func OptionTUI(args OptionTUIArgs) error {
	if args.LogFileName != "" {
		closeLogFile, _ := cmdUtils.ConfigureBubbleTeaLogger(args.LogFileName)
		defer closeLogFile()
	}

	m, err := NewModel(args.Scopes, args.SelectedScopeName, args.MinScore, args.DebounceTime, args.InitialInput)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
