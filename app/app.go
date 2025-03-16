package main

import (
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
	rawJsonData    interface{}
	queryEval      string
	queryHist      string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.jsonOutputView.Init(), m.jsonKeyView.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update window size
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Printf("Size msg in main: %#v", msg)
		msgForChild := tea.WindowSizeMsg{Width: msg.Width / 2, Height: msg.Height}
		m.jsonKeyView, _ = m.jsonKeyView.Update(msgForChild)
		m.jsonOutputView, _ = m.jsonOutputView.Update(msgForChild)
		return m, nil
	case queryview.SelectionChangedMsg:
		log.Printf("Receive selection changed msg %v", msg)
		selectedValue := string(msg)
		m = m.UpdateJsonOutputViewByQuery(selectedValue)
		return m, nil
	case tea.KeyMsg:
		log.Printf("app: KeyMsg: %#v", msg)
		switch msg.Type {
		case tea.KeyEnter:
			query := m.jsonKeyView.CurrentQuery()
			m = m.UpdateJsonOutputViewByQuery(query)
			return m, nil
		case tea.KeyCtrlS:
			// Reload based on current query
			m.queryHist += m.queryEval + "|"
			return m, readOnlyFileExecTeaCmd(m.jsonOutputView.GetJsonData())
		}
	case editorFinishedMsg:
		if msg.err != nil {
			panic(msg.err)
		}
		defer os.Remove(msg.tempJsonPath)
		jsonData, err := LoadJsonFile(msg.tempJsonPath)
		if err != nil {
			panic(err)
		}
		m.SetJsonData(jsonData)
		return m, m.Init()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.jsonKeyView, cmd = m.jsonKeyView.Update(msg)
	cmds = append(cmds, cmd)
	m.jsonOutputView, cmd = m.jsonOutputView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func readOnlyFileExecTeaCmd(jsonData interface{}) tea.Cmd {
	// Save jsonFile and view in editor
	jsonPath, err := getTempJsonPath()
	if err != nil {
		panic(err)
	}
	err = SaveJsonToFile(jsonPath, jsonData)
	if err != nil {
		panic(err)
	}
	cmd := editorFileExecCmd(jsonPath, true)

	// Open editor
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorFinishedMsg{tempJsonPath: jsonPath, err: err}
	})
}

func (m model) UpdateJsonOutputViewByQuery(query string) model {
	var jsonData interface{}
	m.queryEval, jsonData = RobustQueryJsonData(query, m.rawJsonData)
	err := m.jsonOutputView.SetJsonData(jsonData)
	if err != nil {
		panic(err)
	}
	m.jsonKeyView.SetComment(m.queryHist + m.queryEval)
	return m
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

func (m *model) SetJsonData(jsonData interface{}) {
	// Set rawJsonData
	m.rawJsonData = jsonData

	// Query view
	keys, err := GetUnnestedKeys(m.rawJsonData)
	if err != nil {
		panic(err)
	}
	// Add `.` to query everything
	keys = append([]string{"."}, keys...)
	engine := KeySearchEngine{
		keys: keys,
	}
	m.jsonKeyView.SetEngine(engine)

	// Output view
	err = m.jsonOutputView.SetJsonData(m.rawJsonData)
	if err != nil {
		panic(err)
	}
}
