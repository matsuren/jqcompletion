package main

import (
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matsuren/jqcompletion/jsonview"
)

const content = `{
  "name": "example",
  "version": "1.0.0",
  "description": "A sample JSON object",
  "longDescription": "Long description, Long description,Long description,Long description,Long description,",
  "properties": {
    "active": true,
    "count": 42
  },
  "name": "example",
  "version": "1.0.0",
  "description": "A sample JSON object",
  "properties": {
    "active": true,
    "count": 42
  },
    "items": [
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "banana", "price": 1000, "count": 4 }
  ],
}
`

func initModel() model {
	component := jsonview.New(80, 20)
	component.SetContent(content)
	return model{component: component}
}

type model struct {
	component jsonview.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.component, cmd = m.component.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.component.View()
}

func main() {
	if os.Getenv("DEBUGLOG") != "" {
		f, err := tea.LogToFile("debug.log", "jsonview")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}
	model := initModel()
	if _, err := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run(); err != nil {
		fmt.Println(err)
	}
}
