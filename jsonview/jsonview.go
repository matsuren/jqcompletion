package jsonview

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
}

var borderStyle = lipgloss.NewStyle().
	PaddingRight(2).
	Width(80).
	Border(lipgloss.NormalBorder())

func New(width, height int) Model {
	v := viewport.New(width, height)
	return Model{
		viewport: v,
	}
}

func (m *Model) SetContent(content string) {
	maxLength := 10000
	log.Printf("len(content) = %v, maxLength = %v", len(content), maxLength)
	if len(content) < maxLength {
		// Too slow for large file
		m.viewport.SetContent(borderStyle.Render(content))
	} else {
		m.viewport.SetContent(content)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.viewport.View()
}
