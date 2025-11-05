package UI

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	components "github.com/pclubiitk/dbcli/Components"
	"github.com/pclubiitk/dbcli/DB"
)

//TODO::using m.Source and m.Dest to
// show tables and add selection rule


func UpdateSelection(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.Step {
	case StepSelectSourceTable:
		m.SourceTableList, cmd = m.SourceTableList.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if item, ok := m.SourceTableList.SelectedItem().(components.ListItem); ok {
					m.SelectedSourceTbl = string(item)
					var query string
					var coloumns []string
					var err error
					if m.SourceCred["dbVendor"] == "oracle" {
						query = "SELECT column_name FROM all_tab_columns WHERE owner = :1 AND table_name = :2 ORDER BY column_id"
						coloumns, err = DB.QueryStrings(m.Source, query,
							strings.ToUpper(m.SourceCred["user"]),
							strings.ToUpper(m.SelectedSourceTbl))
					} else if m.SourceCred["dbVendor"] == "mysql" {
						query = "SELECT column_name FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position"
						coloumns, err = DB.QueryStrings(m.Source, query,
							strings.ToLower(m.SourceCred["user"]),
							strings.ToLower(m.SelectedSourceTbl))
					}


					
					if err != nil {
						m.ErrMsg = fmt.Sprintf("Failed to list dest tables: %v", err)
						m.SourceColumns = []string{}
					} else {
						m.SourceColumns = coloumns
						m.SourceColumnList = components.CreateList(m.SourceColumns, "📦 Source Columns",m.Width,m.Height)
					}
					m.Step = StepSelectSourceColumns
				}
			case "ctrl+c":
				return m, tea.Quit
			}
		}

	case StepSelectSourceColumns:
		m.SourceColumnList, cmd = m.SourceColumnList.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				// Toggle selected columns
				if item, ok := m.SourceColumnList.SelectedItem().(components.ListItem); ok {
					name := string(item)
					if Contains(m.SelectedSourceCols, name) {
						m.SelectedSourceCols = Remove(m.SelectedSourceCols, name)
					} else {
						m.SelectedSourceCols = append(m.SelectedSourceCols, name)
					}
				}
			case "enter":
				m.Step = StepSelectDestTable
			}
		}

	case StepSelectDestTable:
		m.DestTableList, cmd = m.DestTableList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if item, ok := m.DestTableList.SelectedItem().(components.ListItem); ok {
					m.SelectedDestTbl = string(item)
					var query string
					var coloumns []string
					var err error
					if m.SourceCred["dbVendor"] == "oracle" {
						query = "SELECT column_name FROM all_tab_columns WHERE owner = :1 AND table_name = :2 ORDER BY column_id"
						coloumns, err = DB.QueryStrings(m.Source, query,
							strings.ToUpper(m.SourceCred["user"]),
							strings.ToUpper(m.SelectedSourceTbl))
					} else if m.SourceCred["dbVendor"] == "mysql" {
						query = "SELECT column_name FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position"
						coloumns, err = DB.QueryStrings(m.Source, query,
							strings.ToLower(m.SourceCred["user"]),
							strings.ToLower(m.SelectedSourceTbl))
					}

					if err != nil {
						m.ErrMsg = fmt.Sprintf("Failed to list dest tables: %v", err)
						m.DestColumns = []string{}
					} else {
						m.DestColumns = coloumns
						m.DestColumnList = components.CreateList(m.DestColumns, "📦 Destination Columns",m.Width,m.Height)
					}
					m.Step = StepSelectDestColumns
				}
			}
		}

	case StepSelectDestColumns:
		m.DestColumnList, cmd = m.DestColumnList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ":
				if item, ok := m.DestColumnList.SelectedItem().(components.ListItem); ok {
					name := string(item)
					if Contains(m.SelectedDestCols, name) {
						m.SelectedDestCols = Remove(m.SelectedDestCols, name)
					} else {
						m.SelectedDestCols = append(m.SelectedDestCols, name)
					}
				}
			case "enter":
				m.Step = StepMapping
			}
		}
	}

	return m, cmd
}


func ViewSelection(m Model) string {
	var sb strings.Builder
	title := titleStyle.Render("🧭 Database Selection")
	sb.WriteString("%v")

	switch m.Step {
	case StepSelectSourceTable:
		return lipgloss.Place(m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				m.SourceTableList.View(),
				"",
				infoStyle.Render("\n↑/↓ to navigate, Enter ↵ to select table"),
			),
		)

	case StepSelectSourceColumns:
		return lipgloss.Place(m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				m.SourceColumnList.View(),
				infoStyle.Render("\n↑/↓ to navigate, Space to toggle, Enter ↵ to confirm"),
			),
		)

	case StepSelectDestTable:
		return lipgloss.Place(m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				m.DestTableList.View(),
				infoStyle.Render("\n↑/↓ to navigate, Enter ↵ to select table"),
			),
		)

	case StepSelectDestColumns:
		return lipgloss.Place(m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				m.DestColumnList.View(),
				infoStyle.Render("\n↑/↓ to navigate, Space to toggle, Enter ↵ to continue"),
			),
		)
	}

	return "Unknown selection step"
}
