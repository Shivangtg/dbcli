package UI

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func UpdateDumpOption(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Always update the text input for dump path
	m.DumpPathInp, cmd = m.DumpPathInp.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			// Step 1: Ask whether they want a dump
			if m.WantDump == nil {
				val := true
				m.WantDump = &val
				m.DumpPathInp.Focus() // move to input field
				return m, cmd
			}

			// Step 2: If they want dump, read the input path
			if *m.WantDump {
				if m.DumpPathInp.Value() == "" {
					m.DumpPath = "./dump.sql"
				} else {
					m.DumpPath = m.DumpPathInp.Value()
				}


				if err := createDumpFile(m.DumpPath, m); err != nil {
					m.ErrMsg = fmt.Sprintf("❌ Failed to create dump: %v", err)
				} else {
					m.ErrMsg = fmt.Sprintf("✅ Dump successfully saved to: %s", m.DumpPath)
				}

			}

			// Step 3: Move to next step
			m.Step++
			return m, cmd

		case "n":
			if m.WantDump == nil {
				val := false
				m.WantDump = &val
				m.Step++
				return m, nil
			}
		}
	}

	return m, cmd
}

func ViewDumpOption(m Model) string {
	title := lipgloss.NewStyle().Bold(true).Underline(true).Render("💾 Dump Option")

	if m.WantDump == nil {
		return fmt.Sprintf(
			"%s\n\nDo you want to dump the data to a file?\n(Press ENTER for Yes, or 'n' for No)",
			title,
		)
	}

	if *m.WantDump {
		info := ""
		if m.ErrMsg != "" {
			info = "\n\n" + m.ErrMsg
		}

		return fmt.Sprintf(
			"%s\n\nEnter dump file path (default: ./dump.sql):\n\n%s\n\nPress ENTER to confirm.%s",
			title,
			m.DumpPathInp.View(),
			info,
		)
	}

	return fmt.Sprintf(
		"%s\n\nSkipping dump creation.\nPress ENTER to continue.",
		title,
	)
}
