package jsonview

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

const content = `{
  "name": "example",
  "version": "1.0.0",
  "description": "A sample JSON object",
  "longDescription": "Long description, Long description,Long description,Long description,Long description,",
  "properties": {
    "active": true,
    "count": 42
  }
}
`

func TestMain(m *testing.M) {
	if os.Getenv("DEBUGLOG") != "" {
		f, err := tea.LogToFile("debug.log", "jsonview")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer func() { _ = f.Close() }()
	} else {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}

	m.Run()
}

func initModel() model {
	component := New(80, 20)
	err := component.SetJsonString(content)
	if err != nil {
		panic(err)
	}

	return model{component: component}
}

type model struct {
	component Model
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

func TestUI(t *testing.T) {
	m := initModel()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(300, 100),
	)
	t.Cleanup(func() {
		if err := tm.Quit(); err != nil {
			t.Fatal(err)
		}
	})
	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	// TODO: Check output if it helps
}
