package template

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	quitting bool
	spinner  *spinner.Model
	error    error

	children []tea.Model
}

func NewTemplateModel(s *spinner.Model, children ...tea.Model) tea.Model {
	return model{
		spinner:  s,
		children: children,
	}
}

func (m model) Init() tea.Cmd {
	initCmds := make([]tea.Cmd, len(m.children)+1)
	for i, childModel := range m.children {
		initCmds[i+1] = childModel.Init()
	}

	initCmds[0] = m.spinner.Tick

	return tea.Batch(
		initCmds...,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		s2, cmd := m.spinner.Update(msg)
		*m.spinner = s2
		return m, cmd
	}

	for i, childModel := range m.children {
		m2, cmd := childModel.Update(msg)
		m.children[i] = m2
		if cmd != nil {
			return m, cmd
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	buff := ""

	for _, childModel := range m.children {
		buff += childModel.View()
	}

	return buff
}
