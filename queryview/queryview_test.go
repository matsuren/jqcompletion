package queryview

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
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
	} else {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}

	m.Run()
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

func initModel() model {
	component := New()
	component.SetItems(fakeData)
	return model{component: component}
}

func TestUISelectByCtrlNCtrlP(t *testing.T) {
	m := initModel()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(70, 30),
	)
	t.Cleanup(func() {
		if err := tm.Quit(); err != nil {
			t.Fatal(err)
		}
	})
	if m.component.SelectedValue() != fakeData[0] {
		t.Errorf("Initial selected value is different %s vs %s",
			m.component.SelectedValue(), fakeData[0])
	}
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlN,
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlN,
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlP,
	})
	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	if tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).(model).component.SelectedValue() != fakeData[1] {
		t.Errorf("Final selected value is different %s vs %s",
			m.component.SelectedValue(), fakeData[1])
	}
}

func TestUITabCompletion(t *testing.T) {
	m := initModel()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(70, 30),
	)
	t.Cleanup(func() {
		if err := tm.Quit(); err != nil {
			t.Fatal(err)
		}
	})
	if m.component.SelectedValue() != fakeData[0] {
		t.Errorf("Initial selected value is different %s vs %s",
			m.component.SelectedValue(), fakeData[0])
	}
	tm.Send(tea.KeyMsg{
		Type: tea.KeyCtrlN,
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyTab,
	})
	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	finalModel := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).(model)
	if finalModel.component.queryInput.Value() != fakeData[1] {
		t.Errorf("actual: %s vs expected: %s",
			finalModel.component.queryInput.Value(), fakeData[1])
	}
}

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
	component := New()
	engine := MockEngine{}
	component.SetEngine(engine)
	return model{component: component}
}

func TestUIDebounceQuery(t *testing.T) {
	m := initModelWithMockEngine()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(70, 30),
	)
	t.Cleanup(func() {
		if err := tm.Quit(); err != nil {
			t.Fatal(err)
		}
	})
	tm.Type("This is test query")
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	time.Sleep(time.Millisecond * 500) // Need to wait more than debounceDuration

	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	finalModel := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).(model)
	list := finalModel.component.list
	if len(list.Items()) < 10 {
		t.Errorf("list is not updated: %#v", list.Items())
	}
}
