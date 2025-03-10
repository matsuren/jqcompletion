package queryview

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	queryInput    textinput.Model
	candidateList []string
	selected      int
	currentQuery  string
	engine        Engine
	comment       string
	styleForInput lipgloss.Style
	styleForList  lipgloss.Style
}

type (
	queryRequestMsg  string
	queryResponseMsg []string
)

const (
	debounceDuration = 100 * time.Millisecond
	// TODO: bubbles.list might be better
	maxNumberOfElements = 10
	uiWidth             = 50
)

type Engine interface {
	Query(string) []string
}

func (m *Model) SetEngine(engine Engine) {
	m.engine = engine
}

func requestDebouncedQuery(query string) tea.Cmd {
	return tea.Tick(debounceDuration, func(_ time.Time) tea.Msg {
		return queryRequestMsg(query)
	})
}

func performQuery(query string, engine Engine) tea.Cmd {
	if engine == nil {
		return nil
	}
	return func() tea.Msg {
		response := engine.Query(query)
		return queryResponseMsg(response)
	}
}

func New() Model {
	ti := textinput.New()
	ti.Focus()

	return Model{
		queryInput:    ti,
		candidateList: []string{},
		selected:      0,

		styleForInput: lipgloss.NewStyle().
			Padding(0, 0).
			Width(uiWidth).
			Border(lipgloss.NormalBorder()),
		styleForList: lipgloss.NewStyle().
			Padding(0, 0).
			Width(uiWidth).
			Border(lipgloss.NormalBorder()),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, performQuery("", m.engine))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Save original to detect changes
	originalQueryInputValue := m.queryInput.Value()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Printf("WindowSize: %#v", msg)
		m.styleForInput = m.styleForInput.Width(msg.Width)
		x, _ := m.styleForInput.GetFrameSize()
		m.queryInput.Width = msg.Width - x

		h := lipgloss.Height(m.queryInputView())
		m.styleForList = m.styleForList.Width(msg.Width).Height(msg.Height - h)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlP, tea.KeyUp:
			m.selected = max(m.selected-1, 0)
			return m, nil
		case tea.KeyCtrlN, tea.KeyDown:
			m.selected = min(m.selected+1, len(m.candidateList)-1)
			return m, nil

		case tea.KeyTab:
			log.Printf("Tab: %#v", msg)
			selectedValue := m.SelectedValue()
			if selectedValue != "" {
				m.queryInput.SetValue(selectedValue)
				m.queryInput.CursorEnd()
			}
		}
	case queryRequestMsg:
		if m.currentQuery == string(msg) {
			return m, performQuery(m.currentQuery, m.engine)
		}

	case queryResponseMsg:
		m.candidateList = []string(msg)
		m.selected = 0
		return m, nil
	}
	var cmd tea.Cmd
	m.queryInput, cmd = m.queryInput.Update(msg)
	m.currentQuery = m.queryInput.Value()
	if originalQueryInputValue != m.currentQuery {
		log.Printf("old value: %s vs new value: %s", originalQueryInputValue, m.queryInput.Value())
		return m, tea.Batch(cmd, requestDebouncedQuery(m.currentQuery))
	}
	return m, cmd
}

func (m Model) queryInputView() string {
	return m.styleForInput.Render(m.queryInput.View() + "\n" + m.comment)
}

func (m Model) candidateListView() string {
	s := "Candidate List\n"

	for i, item := range m.candidateList {
		cursor := " "
		if m.selected == i {
			cursor = "|"
		}
		s += fmt.Sprintf("%s %s", cursor, item)
		s += "\n"
		if i > maxNumberOfElements {
			s += "and more...\n"
			break
		}
	}
	return m.styleForList.Render(s)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.queryInputView(), m.candidateListView())
}

func (m Model) SelectedValue() string {
	if m.selected < len(m.candidateList) {
		return m.candidateList[m.selected]
	}
	return ""
}

func (m Model) CurrentQuery() string {
	return m.queryInput.Value()
}

func (m *Model) SetQueryInput(query string) {
	m.queryInput.SetValue(query)
}

func (m *Model) SetComment(comment string) {
	m.comment = comment
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
