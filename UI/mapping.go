package UI

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	components "github.com/pclubiitk/dbcli/Components"
)

//TODO::Put logic to map selected coloumns of both source
//and destination database



func UpdateMapping(m Model, msg tea.Msg) (Model,tea.Cmd) {
	var cmd tea.Cmd
	m.DestColumnList, cmd = m.DestColumnList.Update(msg)

	if m.ColumnMapping == nil {
		m.ColumnMapping = make(map[string]string)
	}
	
	if m.CurrentMapIdx < len(m.SourceColumns) {
        m.ColumnBeingMapped = m.SourceColumns[m.CurrentMapIdx]
    }

	switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            // Get selected destination column
            if item, ok := m.DestColumnList.SelectedItem().(components.ListItem); ok {
                destCol := item.Title()
                m.ColumnMapping[m.ColumnBeingMapped] = destCol
                m.CurrentMapIdx++
                
                // If finished all mappings, go to next step
                if m.CurrentMapIdx >= len(m.SourceColumns) {
                    m.Step++
                } else {
                    m.ColumnBeingMapped = m.SourceColumns[m.CurrentMapIdx]
                }
            }
        }
    }
	
	return m, cmd
}

func ViewMapping(m Model) string {
    title := titleStyle.Render(fmt.Sprintf("🔗 Mapping source column: %s", m.ColumnBeingMapped))

    progress := infoStyle.Render(fmt.Sprintf(
        "Mapped %d / %d columns",
        len(m.ColumnMapping),
        len(m.SourceColumns),
    ))

    return lipgloss.Place(
        m.Width,
        m.Height,
        lipgloss.Center,
        lipgloss.Center,
        lipgloss.JoinVertical(
            lipgloss.Left,
            title,
            progress,
            "",
            m.DestColumnList.View(),
            "",
            quitStyle.Render("Press Enter to confirm mapping."),
        ),
    )
}
