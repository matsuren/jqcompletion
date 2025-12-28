package main

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

func TestAppUISampleData(t *testing.T) {
	m := initializeModel()
	jsonData, err := LoadJsonFile("../data/sample.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	m.SetJsonData(jsonData)
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(100, 40),
	)
	time.Sleep(time.Millisecond * 100) // Wait for model initialization
	tm.Type("rel")
	time.Sleep(time.Millisecond * 400) // Wait for debounceDuration
	tm.Send(tea.KeyMsg{
		Type: tea.KeyTab,
	})
	time.Sleep(time.Millisecond * 100)

	// Quit app and check output
	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	finalModel := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).(model)

	// Check json output view
	want := `[
  "Info",
  "Debug"
]`
	got := jsonDataToStrings(finalModel.jsonOutputView.GetJsonData())
	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}

	// Check json key view
	want = ".results[].level"
	got = finalModel.jsonKeyView.CurrentQuery()
	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
	numberOfItems := len(finalModel.jsonKeyView.GetItems())
	if numberOfItems != 1 {
		t.Errorf("Got: %v. Number of items in jsonKeyView is wrong.", numberOfItems)
	}
	got = finalModel.jsonKeyView.GetItems()[0]
	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestAppUISampleDataTypeEnter(t *testing.T) {
	m := initializeModel()
	jsonData, err := LoadJsonFile("../data/sample.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	m.SetJsonData(jsonData)
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(100, 40),
	)
	time.Sleep(time.Millisecond * 100) // Wait for model initialization
	tm.Type(".results[0]|keys")
	time.Sleep(time.Millisecond * 400)
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// Quit app and check output
	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	finalModel := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second)).(model)

	// Check json output view
	want := `[
  "level",
  "message"
]`
	got := jsonDataToStrings(finalModel.jsonOutputView.GetJsonData())
	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}
