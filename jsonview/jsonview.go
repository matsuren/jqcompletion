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
	content  string
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

func (m *Model) SetContent(content string) {
	m.viewport.SetContent(content)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Width(msg.Width).Height(msg.Height)
		x, _ := m.style.GetFrameSize()
		m.viewport.Width = msg.Width - x
		m.viewport.Height = msg.Height
		log.Printf("JsonView msg: %#v", msg)
		return m, nil

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
	return m.style.Render(m.viewport.View())
}
