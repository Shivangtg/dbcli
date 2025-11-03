package UI

import "github.com/charmbracelet/lipgloss"

// Define reusable styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00BFFF"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			MarginTop(1)

	msgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#632496ff"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F")).
			Bold(true)

	quitStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			MarginTop(1)
)
