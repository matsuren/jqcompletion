package main

import (
	"encoding/json"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matsuren/jqcompletion/jsonview"
	"github.com/matsuren/jqcompletion/queryview"
)

type model struct {
	queryView queryview.Model
	jsonView  jsonview.Model
	jsonData  interface{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Start")
	var cmds []tea.Cmd

	// Update jsonView
	jsonViewModel, cmd := m.jsonView.Update(msg)
	if updatedView, ok := jsonViewModel.(jsonview.Model); ok {
		m.jsonView = updatedView
	} else {
		panic("Wrong type")
	}
	cmds = append(cmds, cmd)

	// Update queryView
	queryViewModel, cmd := m.queryView.Update(msg)
	log.Printf("query %#v", cmd)
	if updatedView, ok := queryViewModel.(queryview.Model); ok {
		m.queryView = updatedView
	} else {
		panic("Wrong type")
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.queryView.View(), m.jsonView.View())
}

func initializeModel() model {
	qv := queryview.New()
	jv := jsonview.New(80, 10)
	return model{
		queryView: qv,
		jsonView:  jv,
	}
}

func (m *model) LoadJsonFile(jsonPath string) {
	log.Println("Loading ", jsonPath)
	// Read the JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	// Parse the JSON
	log.Println("Parsing ", jsonPath)
	err = json.Unmarshal(jsonData, &m.jsonData)
	if err != nil {
		panic(err)
	}
	log.Println("Done LoadJsonFile")
}

func (m *model) SetJsonDataInView(jsonData interface{}) {
	log.Println("Start json.MarshalIndent")
	resultBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println("Start SetContent")
	m.jsonView.SetContent(string(resultBytes))
	log.Println("Done SetJsonDataInView")
}
