package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestMain(m *testing.M) {
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
	tm.Type("re")
	time.Sleep(time.Millisecond * 400) // Wait for debounceDuration
    tm.Send(tea.KeyMsg{
        Type: tea.KeyTab,
    })

	// Quit app and check output
	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}
