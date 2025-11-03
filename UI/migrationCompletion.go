package UI

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type migrationMsg struct {
	Progress int
}
type migrationDoneMsg struct {
	Err error
}

// UpdateMigrationCompletion handles UI updates during migration.
func UpdateMigrationCompletion(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.MigrationInProgress {
				return m, nil
			}
			m.MigrationInProgress = true
			m.StatusMsg = "Starting migration..."

			progressChan := make(chan int)
			doneChan := make(chan error)

			// Run migration concurrently
			go func() {
				MigrateData(m, progressChan, doneChan)
			}()

			// Wait for both progress and done messages
			return m, tea.Batch(
				waitForProgress(progressChan),
				waitForDone(doneChan),
			)

		case "q":
			return m, tea.Quit
		}

	case migrationMsg:
		m.Progress.SetPercent(float64(msg.Progress) / 100.0)
		m.StatusMsg = fmt.Sprintf("Migrating... %d%%", msg.Progress)
		// Keep listening for more progress
		return m, waitForProgress(m.ProgressChan)

	case migrationDoneMsg:
		m.MigrationInProgress = false
		if msg.Err != nil {
			m.ErrMsg = fmt.Sprintf("❌ Migration failed: %v", msg.Err)
			m.StatusMsg = ""
		} else {
			m.Progress.SetPercent(1.0)
			m.StatusMsg = "✅ Migration completed successfully!"
		}
		return m, nil
	}

	return m, nil
}

func waitForProgress(ch <-chan int) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return nil
		}
		return migrationMsg{Progress: p}
	}
}

func waitForDone(ch <-chan error) tea.Cmd {
	return func() tea.Msg {
		err := <-ch
		return migrationDoneMsg{Err: err}
	}
}

func ViewMigrationCompletion(m Model) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("82")).
		Render("🚀 Data Migration")

	progressView := m.Progress.View()

	status := ""
	if m.StatusMsg != "" {
		status = lipgloss.NewStyle().Render(m.StatusMsg)
	}
	err := ""
	if m.ErrMsg != "" {
		err = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.ErrMsg)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		progressView,
		"",
		status,
		err,
		"",
		lipgloss.NewStyle().Faint(true).Render("Press Q to quit."),
	)
}
