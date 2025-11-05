package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListItem defines each row in the list
type ListItem string

func (i ListItem) Title() string       { return string(i) }
func (i ListItem) Description() string { return "" }
func (i ListItem) FilterValue() string { return string(i) }


//
// Custom Delegate (for toggled selection coloring)
//

type CustomDelegate struct {
	selectedItems map[string]bool
}

func NewCustomDelegate() *CustomDelegate {
	return &CustomDelegate{selectedItems: make(map[string]bool)}
}

func (d *CustomDelegate) SetSelected(selected []string) {
	d.selectedItems = make(map[string]bool)
	for _, s := range selected {
		d.selectedItems[s] = true
	}
}

func (d *CustomDelegate) Height() int  { return 1 }
func (d *CustomDelegate) Spacing() int { return 0 }

// ✅ FIXED: use tea.Msg instead of list.Message
func (d *CustomDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d *CustomDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ListItem)
	if !ok {
		return
	}

	name := string(i)
	cursor := "  "
	if index == m.Index() {
		cursor = "➤ "
	}

	// Define styles
	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)

	// Determine which style applies
	style := defaultStyle
	if d.selectedItems[name] {
		style = selectedStyle
	}
	if index == m.Index() {
		style = cursorStyle
	}

	// Add ✅ checkmark for toggled items
	check := ""
	if d.selectedItems[name] {
		check = " ✅"
	}

	_, _ = fmt.Fprintln(w, style.Render(cursor+name+check))
}

// --- List Creation ---

func CreateList(items []string, title string, width int, height int) list.Model {
	li := make([]list.Item, len(items))
	for i, t := range items {
		li[i] = ListItem(t)
	}

	delegate := NewCustomDelegate()

	l := list.New(li, delegate, width, height)
	l.Title = title
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	return l
}
