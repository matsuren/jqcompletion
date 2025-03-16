package jsonview

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	style    lipgloss.Style
	jsonData interface{}
}

func New(width, height int) Model {
	v := viewport.New(width, height)
	return Model{
		viewport: v,
		style: lipgloss.NewStyle().
			Width(width).
			Border(lipgloss.NormalBorder()),
	}
}

func (m *Model) SetJsonString(jsonString string) error {
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON in SetJsonString: %w", err)
	}
	err = m.SetJsonData(jsonData)
	if err != nil {
		return fmt.Errorf("error SetJsonData in SetJsonString: %w", err)
	}
	return nil
}

func (m *Model) SetJsonData(jsonData interface{}) error {
	log.Println("Start json.MarshalIndent")
	resultBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent failed: %w", err)
	}
	log.Println("Start SetContent")
	m.viewport.SetContent(string(resultBytes))
	log.Println("Done SetJsonDataInView")

	// Save jsonData for getter
	m.jsonData = jsonData
	return nil
}

func (m Model) GetJsonData() interface{} {
	return m.jsonData
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// FIXME: Somehow margin is needed to avoid hidden top
		margin := 3
		m.style = m.style.Width(msg.Width - margin).Height(msg.Height - margin)
		x, _ := m.style.GetFrameSize()
		m.viewport.Width = m.style.GetWidth() - x
		m.viewport.Height = m.style.GetHeight()
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
