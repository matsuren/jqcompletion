package main

import (
	// "github.com/itchyny/gojq"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app            *tview.Application
	queryInput     *tview.InputField
	evaluatedQuery *tview.TextView
	jsonOutput     *tview.TextView
	suggestions    *tview.List
	layout         *tview.Flex
}

func NewJSONQueryUI() *UI {
	ui := &UI{
		app:            tview.NewApplication(),
		queryInput:     tview.NewInputField(),
		evaluatedQuery: tview.NewTextView(),
		jsonOutput:     tview.NewTextView(),
		suggestions:    tview.NewList(),
		layout:         tview.NewFlex(),
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
	suggestions := []string{"Test0", "Test1", "Test2"}
	for _, v := range suggestions {
		ui.suggestions.AddItem(v, "", 0, nil)
	}

	// Key binding
	ui.queryInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlP:
			current := ui.suggestions.GetCurrentItem()
			if current > 0 {
				ui.suggestions.SetCurrentItem(current - 1)
			}
			return nil
		case event.Key() == tcell.KeyCtrlN:
			current := ui.suggestions.GetCurrentItem()
			if current < ui.suggestions.GetItemCount()-1 {
				ui.suggestions.SetCurrentItem(current + 1)
			}
			return nil
		case event.Key() == tcell.KeyTab:
			if ui.suggestions.GetItemCount() > 0 {
				index := ui.suggestions.GetCurrentItem()
				mainText, _ := ui.suggestions.GetItemText(index)
				currentQuery := ui.queryInput.GetText()
				ui.jsonOutput.SetText(fmt.Sprintf("Debug: %v, %v", mainText, currentQuery))
			}
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
		AddItem(ui.jsonOutput, 0, 1, false)

	ui.layout = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFlex, 0, 1, true).
		AddItem(rightFlex, 0, 1, false)

	return ui
}

func (ui *UI) Run() error {
	ui.app.SetRoot(ui.layout, true)
	return ui.app.Run()
}
