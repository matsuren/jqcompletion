package queryview

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	style := lipgloss.NewStyle().PaddingLeft(3)
	frameX, _ := style.GetFrameSize()
	str := fmt.Sprintf("%s", i)
	if len(str)+frameX > m.Width() {
		prefix := "…"
		str = prefix + str[len(str)-m.Width()+frameX+len(prefix):]
	}
	if index == m.Index() {
		prefix := "> "
		str = prefix + str
		newPaddingLeft := style.GetPaddingLeft() - len(prefix)
		if newPaddingLeft < 0 {
			panic("Negative padding. Check len(prefix)")
		}
		style = style.PaddingLeft(newPaddingLeft).Foreground(mainColor)
	}
	fmt.Fprint(w, style.Render(str))
}

func newList(width, height int) list.Model {
	items := []list.Item{}
	l := list.New(items, itemDelegate{}, width, height)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(2)
	return l
}
