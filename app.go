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
	jsonKeyView    queryview.Model
	jsonOutputView jsonview.Model
	rawJsonData       interface{}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.jsonOutputView.Init(), m.jsonKeyView.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Start")
	var cmds []tea.Cmd

	// Update jsonView
	jsonViewModel, cmd := m.jsonOutputView.Update(msg)
	if updatedView, ok := jsonViewModel.(jsonview.Model); ok {
		m.jsonOutputView = updatedView
	} else {
		panic("Wrong type")
	}
	cmds = append(cmds, cmd)

	// Update queryView
	queryViewModel, cmd := m.jsonKeyView.Update(msg)
	log.Printf("query %#v", cmd)
	if updatedView, ok := queryViewModel.(queryview.Model); ok {
		m.jsonKeyView = updatedView
	} else {
		panic("Wrong type")
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.jsonKeyView.View(), m.jsonOutputView.View())
}

func initializeModel() model {
	qv := queryview.New()
	jv := jsonview.New(80, 10)
	return model{
		jsonKeyView:    qv,
		jsonOutputView: jv,
	}
}

func initializeModelWithJsonFile(jsonPath string) model {
	m := initializeModel()

	m.LoadJsonFile(jsonPath)
	m.setJsonDataInView(m.rawJsonData)

	keys, err := GetUnnestedKeys(m.rawJsonData)
	if err != nil {
		panic(err)
	}
	engine := KeySearchEngine{
		keys: keys,
	}
	m.jsonKeyView.SetEngine(engine)
	return m
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
	err = json.Unmarshal(jsonData, &m.rawJsonData)
	if err != nil {
		panic(err)
	}
	log.Println("Done LoadJsonFile")
}

func (m *model) setJsonDataInView(jsonData interface{}) {
	log.Println("Start json.MarshalIndent")
	resultBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println("Start SetContent")
	m.jsonOutputView.SetContent(string(resultBytes))
	log.Println("Done SetJsonDataInView")
}
