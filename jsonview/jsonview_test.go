// TODO: enable go:build debug
package jsonview

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

const content = `
{
  "name": "example",
  "version": "1.0.0",
  "description": "A sample JSON object",
  "longDescription": "Long description, Long description,Long description,Long description,Long description,",
  "properties": {
    "active": true,
    "count": 42
  },
  "name": "example",
  "version": "1.0.0",
  "description": "A sample JSON object",
  "properties": {
    "active": true,
    "count": 42
  },
    "items": [
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "apple", "price": 110, "count": 3 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "orange", "price": 120, "count": 2 },
    { "name": "orange2", "price": null, "count": 2 },
    { "name": "banana", "price": 1000, "count": 4 }
  ],
}
`

func TestUIJsonView(t *testing.T) {
	model := New(100, 20)
	model.SetContent(content)
	if _, err := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run(); err != nil {
		t.Error(err)
	}
}

