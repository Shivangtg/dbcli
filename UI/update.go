package UI

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd      // command returned from sub-updates
	var cmds []tea.Cmd   // collect multiple commands if needed

	switch m.Step {
	case StepSourceCred, StepDestCred:
		m, cmd = UpdateDBCred(m, msg)
		cmds = append(cmds, cmd)

	case StepSelectSourceTable, StepSelectSourceColumns, StepSelectDestTable, StepSelectDestColumns:
		m, cmd = UpdateSelection(m, msg)
		cmds = append(cmds, cmd)

	case StepMapping:
		m, cmd = UpdateMapping(m, msg)
		cmds = append(cmds, cmd)

	case StepDumpOption:
		m, cmd = UpdateDumpOption(m, msg)
		cmds = append(cmds, cmd)

	case StepMigrationConfirm:
		m, cmd = UpdateMigrationCompletion(m, msg)
		cmds = append(cmds, cmd)
	}

	// Global key handling
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Return the model and batch all commands together
	return m, tea.Batch(cmds...)
}

