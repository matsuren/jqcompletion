package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app            *tview.Application
	queryInput     *tview.InputField
	evaluatedQuery *tview.TextView
	jsonOutput     *tview.TextView
	debugLog       *tview.TextView
	suggestions    *tview.List
	layout         *tview.Flex
	jsonData       interface{}
}

func NewJSONQueryUI() *UI {
	ui := &UI{
		app:            tview.NewApplication(),
		queryInput:     tview.NewInputField(),
		evaluatedQuery: tview.NewTextView(),
		jsonOutput:     tview.NewTextView(),
		suggestions:    tview.NewList(),
		debugLog:       tview.NewTextView(),
		layout:         tview.NewFlex(),
		jsonData:       nil,
	}

	ui.queryInput.SetTitle("Query")
	ui.queryInput.SetBorder(true)
	ui.queryInput.SetFieldWidth(0)
	ui.queryInput.SetLabel("jq> ")

	ui.evaluatedQuery.SetTitle("evaluated query")
	ui.evaluatedQuery.SetBorder(true)

	ui.jsonOutput.SetTitle("Result")
	ui.jsonOutput.SetBorder(true)

	ui.suggestions.SetTitle("Suggestions")
	ui.suggestions.ShowSecondaryText(false)
	ui.suggestions.SetBorder(true)

	ui.debugLog.SetTitle("debug")
	ui.debugLog.SetBorder(true)

	// Key binding
	ui.queryInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlP:
			current := ui.suggestions.GetCurrentItem()
			if current > 0 {
				ui.suggestions.SetCurrentItem(current - 1)
				ui.UpdateQueryJsonDataBySelection()
			}
			return nil
		case event.Key() == tcell.KeyCtrlN:
			current := ui.suggestions.GetCurrentItem()
			if current < ui.suggestions.GetItemCount()-1 {
				ui.suggestions.SetCurrentItem(current + 1)
				ui.UpdateQueryJsonDataBySelection()
			}
			return nil
		case event.Key() == tcell.KeyTab:
			if ui.suggestions.GetItemCount() > 0 {
				index := ui.suggestions.GetCurrentItem()
				query, _ := ui.suggestions.GetItemText(index)
				currentQuery := ui.queryInput.GetText()
				ui.queryInput.SetText(query)
				ui.QueryJsonDataAndShow(query)
				ui.SetDebugText(fmt.Sprintf("Debug: %v, %v", query, currentQuery))
			}
			return nil
		case event.Key() == tcell.KeyEnter:
			query := ui.queryInput.GetText()
			ui.QueryJsonDataAndShow(query)
			return nil
		}
		return event
	})

	// Create layout
	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.queryInput, 3, 0, true).
		AddItem(ui.suggestions, 0, 1, false)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.evaluatedQuery, 3, 0, false).
		AddItem(ui.jsonOutput, 0, 1, false).
		AddItem(ui.debugLog, 3, 0, false)

	ui.layout = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFlex, 0, 1, true).
		AddItem(rightFlex, 0, 1, false)

	return ui
}

func (ui *UI) UpdateQueryJsonDataBySelection() {
	index := ui.suggestions.GetCurrentItem()
	query, _ := ui.suggestions.GetItemText(index)
	ui.QueryJsonDataAndShow(query)
}

func (ui *UI) SetDebugText(text string) {
	ui.debugLog.SetText(text)
}

func (ui *UI) SetOutputJsonData(jsonData interface{}) {
	resultBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		ui.SetDebugText(fmt.Sprintf("Error: %v", err))
	}
	ui.jsonOutput.SetText(string(resultBytes))
}

func (ui *UI) QueryJsonDataAndShow(query string) {
	evalQuery, result := RobustQueryJsonData(query, ui.jsonData)
	if len(evalQuery) > 0 {
		ui.evaluatedQuery.SetText(evalQuery)
		ui.SetOutputJsonData(result)
	} else if msg, ok := result.(string); ok {
		ui.SetDebugText(msg)
	}
}

func (ui *UI) LoadJsonFile(jsonPath string) error {
	fmt.Println("Loading ", jsonPath)
	// Read the JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("Error reading file: %v\n", err)
	}

	// Parse the JSON
	err = json.Unmarshal(jsonData, &ui.jsonData)
	if err != nil {
		return fmt.Errorf("Error parsing JSON: %v\n", err)
	}
	return nil
}

func (ui *UI) Run() error {
	if ui.jsonData == nil {
		return fmt.Errorf("Error: LoadJsonFile before Run")
	}

	ui.SetOutputJsonData(ui.jsonData)

	unnestedKeys, err := GetUnnestedKeys(ui.jsonData)
	if err != nil {
		return fmt.Errorf("Error parsing JSON: %v\n", err)
	}
	for _, key := range unnestedKeys {
		ui.suggestions.AddItem(key, "", 0, nil)
	}

	return ui.app.SetRoot(ui.layout, true).EnableMouse(true).Run()
}
