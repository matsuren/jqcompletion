package queryview

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	queryInput    textinput.Model
	list          list.Model
	engine        Engine
	comment       string
	styleForInput lipgloss.Style
	styleForList  lipgloss.Style
}

var mainColor = lipgloss.Color("63")

type (
	SelectionChangedMsg string
	queryRequestMsg     string
	queryResponseMsg    []string
)

const (
	debounceDuration = 100 * time.Millisecond
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
	ti.PromptStyle = ti.PromptStyle.Foreground(mainColor)
	ti.Prompt = "jq: "

	l := newList(20, 30)

	return Model{
		queryInput: ti,
		list:       l,
		styleForInput: lipgloss.NewStyle().
			Padding(0, 0).
			Border(lipgloss.NormalBorder()),
		styleForList: lipgloss.NewStyle().
			Padding(0, 0).
			Border(lipgloss.NormalBorder()),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, performQuery("", m.engine))
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
		m.list.SetWidth(msg.Width - x)
		m.list.SetHeight(m.styleForList.GetHeight())
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlP, tea.KeyUp:
			m.list.CursorUp()
			return m, m.selectionChanged()
		case tea.KeyCtrlN, tea.KeyDown:
			m.list.CursorDown()
			return m, m.selectionChanged()

		case tea.KeyTab:
			log.Printf("Tab: %#v", msg)
			selectedValue, ok := m.list.SelectedItem().(item)
			if ok {
				m.queryInput.SetValue(string(selectedValue))
				m.queryInput.CursorEnd()
			}
			return m, nil
		}
	case queryRequestMsg:
		if m.queryInput.Value() == string(msg) {
			return m, performQuery(m.queryInput.Value(), m.engine)
		}

	case queryResponseMsg:
		m.SetItems([]string(msg))
		return m, m.selectionChanged()
	}
	var cmd tea.Cmd
	m.queryInput, cmd = m.queryInput.Update(msg)
	if originalQueryInputValue != m.queryInput.Value() {
		log.Printf("old value: %s vs new value: %s", originalQueryInputValue, m.queryInput.Value())
		return m, tea.Batch(cmd, requestDebouncedQuery(m.queryInput.Value()))
	}
	return m, cmd
}

func (m Model) queryInputView() string {
	comment := m.comment
	if len(comment) > m.queryInput.Width {
		comment = comment[:m.queryInput.Width] + ".."
	}
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Dark: "#9B9B9B", Light: "#5C5C5C"})
	return m.styleForInput.Render(m.queryInput.View() + "\n" + commentStyle.Render(comment))
}

func (m Model) listView() string {
	return m.styleForList.Render(m.list.View())
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.queryInputView(), m.listView())
}

func (m *Model) SetItems(items []string) {
	listItems := make([]list.Item, 0, len(items))
	for _, listitem := range items {
		listItems = append(listItems, item(listitem))
	}
	m.list.SetItems(listItems)
	m.list.ResetSelected()
}

func (m Model) SelectedValue() string {
	selectedValue, ok := m.list.SelectedItem().(item)
	if ok {
		return string(selectedValue)
	}
	return ""
}

func (m Model) selectionChanged() tea.Cmd {
	log.Printf("Changed to new value: %v", m.SelectedValue())
	return func() tea.Msg {
		return SelectionChangedMsg(m.SelectedValue())
	}
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
