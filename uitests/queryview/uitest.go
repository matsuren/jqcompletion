package main

import (
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matsuren/jqcompletion/queryview"
)

type MockEngine struct{}

func (e MockEngine) Query(query string) []string {
	list := make([]string, 0, len(query))
	list = append(list, "default 1")
	for i := 0; i < len(query); i++ {
		list = append(list, query[:i+1])
	}
	return list
}

func initModelWithMockEngine() model {
	component := queryview.New()
	engine := MockEngine{}
	component.SetEngine(engine)
	return model{component: component}
}

type model struct {
	component queryview.Model
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
		f, err := tea.LogToFile("debug.log", "queryview")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer func() { _ = f.Close() }()
	} else {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}

	model := initModelWithMockEngine()
	if _, err := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	).Run(); err != nil {
		fmt.Println(err)
	}
}
