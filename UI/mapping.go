package UI

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	components "github.com/pclubiitk/dbcli/Components"
)

// UpdateMapping handles interactive mapping between selected source & destination columns
func UpdateMapping(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.MappingDestColumnList, cmd = m.MappingDestColumnList.Update(msg)

	// Initialize mapping storage if nil
	if m.ColumnMapping == nil {
		m.ColumnMapping = make(map[string]string)
	}

	// Make sure we have a consistent list of source columns to map
	sourceCols := m.SelectedSourceCols
	if len(sourceCols) == 0 {
		sourceCols = m.SourceColumns
	}

	// Guard for index bounds
	if m.CurrentMapIdx < len(sourceCols) {
		m.ColumnBeingMapped = sourceCols[m.CurrentMapIdx]
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			if item, ok := m.MappingDestColumnList.SelectedItem().(components.ListItem); ok {
				destCol := item.Title()
				sourceCol := m.ColumnBeingMapped

				// Save mapping
				m.ColumnMapping[sourceCol] = destCol
				m.CurrentMapIdx++

				// Check if all mappings are done
				if m.CurrentMapIdx >= len(sourceCols) {
					m.Step++ // Move to next UI step
                    m.DumpPathInp = &components.Field{}
                    m.DumpPathInp.Blur()
				} else {
					// Move to next source column for mapping
					m.ColumnBeingMapped = sourceCols[m.CurrentMapIdx]
				}
			}

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, cmd
}

// ViewMapping renders the mapping screen
func ViewMapping(m Model) string {
	sourceCols := m.SelectedSourceCols
	if len(sourceCols) == 0 {
		sourceCols = m.SourceColumns
	}

	progress := fmt.Sprintf("Mapped %d / %d columns", len(m.ColumnMapping), len(sourceCols))
	title := titleStyle.Render(fmt.Sprintf("🔗 Mapping source column: %s", m.ColumnBeingMapped))

	var currentMap string
	if dest, ok := m.ColumnMapping[m.ColumnBeingMapped]; ok {
		currentMap = infoStyle.Render(fmt.Sprintf("✔ %s → %s", m.ColumnBeingMapped, dest))
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			infoStyle.Render(progress),
			currentMap,
			"",
			m.MappingDestColumnList.View(),
			"",
			quitStyle.Render("↑/↓ to navigate, Enter ↵ to map, Ctrl+C to quit"),
		),
	)
}
