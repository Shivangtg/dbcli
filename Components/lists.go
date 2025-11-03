package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func CreateList(items []string, title string,width int ,height int) list.Model {
	li := make([]list.Item, len(items))
	for i, t := range items {
		li[i] = ListItem(t)
	}
	l := list.New(li, list.NewDefaultDelegate(),width,height)
	l.Title = title
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	return l
}

type ListItem string

func (i ListItem) Title() string       { return string(i) }
func (i ListItem) Description() string { return "" }
func (i ListItem) FilterValue() string { return string(i) }
