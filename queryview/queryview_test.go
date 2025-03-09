// TODO: enable go:build debug
package queryview

import (
	"fmt"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

var fakeData = []string{"AAA", "BBB", "CCC", "DDD"}

func TestMain(m *testing.M) {
	if os.Getenv("DEBUGLOG") != "" {
		f, err := tea.LogToFile("debug.log", "queryview")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	m.Run()
}

func TestUIOnlyView(t *testing.T) {
	model := New()
	model.candidateList = fakeData
	if _, err := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	).Run(); err != nil {
		t.Error(err)
	}
}

type MockEngine struct{}

func (e MockEngine) Query(query string) []string {
	list := make([]string, 0, len(query))
	for i := 0; i < len(query); i++ {
		list = append(list, query[:i+1])
	}
	return list
}

func TestUIDebounceQuery(t *testing.T) {
	model := New()
	engine := MockEngine{}
	model.SetEngine(engine)
	if _, err := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	).Run(); err != nil {
		t.Error(err)
	}
}
