package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	components "github.com/pclubiitk/dbcli/Components"
	"github.com/pclubiitk/dbcli/UI"
)


func main() {
	model := UI.Model{
		Step:     UI.StepSourceCred,
		CredKeys: []string{"dbVendor", "host", "port", "user", "password", "dbname"},
		IsSource: true,
		Progress: progress.New(progress.WithDefaultGradient()),
		CredInput: &components.Field{},
	}

	p := tea.NewProgram(model,tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

