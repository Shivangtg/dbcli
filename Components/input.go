package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Msg
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

// We need to have a wrapper for our bubbles as they don't currently implement the tea.Model interface
// textinput, textarea

type Field struct {
	textinput textinput.Model
}

func NewField() *Field {
	a := Field{}

	model := textinput.New()
	model.Placeholder = "Your answer here"
	model.Focus()

	a.textinput = model
	return &a
}

func (a *Field) Blink() tea.Msg {
	return textinput.Blink
}

func (a *Field) Init() tea.Cmd {
	return nil
}

func (a *Field) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *Field) View() string {
	return a.textinput.View()
}

func (a *Field) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *Field) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *Field) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *Field) Value() string {
	return a.textinput.Value()
}
