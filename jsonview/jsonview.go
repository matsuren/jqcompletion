package jsonview

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	style    lipgloss.Style
}

func New(width, height int) Model {
	v := viewport.New(width, height)
	return Model{
		viewport: v,
		style: lipgloss.NewStyle().
			PaddingRight(2).
			Width(width).
			Border(lipgloss.NormalBorder()),
	}
}

func (m *Model) SetWidth(width int) {
	m.style = m.style.Width(width -3)
	m.viewport.Width = width
}
func (m *Model) SetHeight(height int) {
	m.viewport.Height = height
}

func (m *Model) SetContent(content string) {
	maxLength := 10000
	log.Printf("len(content) = %v, maxLength = %v", len(content), maxLength)
	if len(content) < maxLength {
		// Too slow for large file
		m.viewport.SetContent(m.style.Render(content))
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
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.viewport.View()
}
