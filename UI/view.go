package UI

import "fmt"

func (m Model) View() string {
	switch m.Step {
	case StepSourceCred, StepDestCred:
		return ViewDBCred(m)
	case StepSelectSourceTable, StepSelectSourceColumns, StepSelectDestTable, StepSelectDestColumns:
		return ViewSelection(m)
	case StepMapping:
		return ViewMapping(m)
	case StepDumpOption:
		return ViewDumpOption(m)
	case StepMigrationConfirm:
		return ViewMigrationCompletion(m)
	default:
		return fmt.Sprintf("Migration complete! Column mappings: %+v\nPress q to quit.", m.ColumnMapping)
	}
}
