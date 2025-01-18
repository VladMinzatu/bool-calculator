package cmd

import (
	"fmt"
	"strings"

	"github.com/VladMinzatu/bool-calculator/evaluation"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type TerminalApp struct{}

func (app TerminalApp) Run() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}

type errMsg error

type model struct {
	input  textinput.Model
	output textarea.Model
	result *evaluation.Result
	err    error
}

func NewModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a boolean expression..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	ta := textarea.New()
	ta.Placeholder = "Output will appear here..."
	ta.ShowLineNumbers = false
	ta.Blur()

	return model{
		input:  ti,
		output: ta,
		result: nil,
		err:    nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	// Validate input as user types
	m.err = m.validateInput()
	if m.err == nil {
		m.output.SetValue(m.result.String())
	} else {
		m.err = fmt.Errorf("* %v", m.err)
		m.output.SetValue("")
	}

	m.output, cmd = m.output.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Enter boolean expression:\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()))
		b.WriteString("\n\n")
	} else {
		b.WriteString("Result:\n")
		b.WriteString(m.output.View())
	}

	b.WriteString("\nPress Esc to quit\n")

	return b.String()
}

func (m *model) validateInput() error {
	result, err := evaluation.Compute(m.input.Value())
	m.result = result
	m.err = err
	return err
}

func (m *model) processValidInput(input string) string {
	return m.result.String()
}
