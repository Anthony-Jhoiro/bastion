package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"os"
)

var (
	color          = termenv.EnvColorProfile().Color
	errorKeyword   = termenv.Style{}.Foreground(color("#E06C75")).Styled
	successKeyword = termenv.Style{}.Foreground(color("#98C379")).Styled
	infoKeyword    = termenv.Style{}.Foreground(color("#61AFEF")).Styled
	heading        = termenv.Style{}.Foreground(color("#61AFEF")).Styled
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var hydraBanner = heading("  ______           _   _             \n  | ___ \\         | | (_)            \n  | |_/ / __ _ ___| |_ _  ___  _ __  \n  | ___ \\/ _` / __| __| |/ _ \\| '_ \\ \n  | |_/ / (_| \\__ \\ |_| | (_) | | | |\n  \\____/ \\__,_|___/\\__|_|\\___/|_| |_|\n                                     ")

type ConnectionStatus = int

const (
	Connected        ConnectionStatus = iota
	ConnectionFailed ConnectionStatus = iota
	Connecting       ConnectionStatus = iota
)

type model struct {
	quitting bool
	spinner  *spinner.Model
	error    error

	children []tea.Model

	vpnModel          VpnModel
	hostSelectorModel HostSelectorModel
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

	buff := fmt.Sprintf("%v\n\n", hydraBanner)

	for _, childModel := range m.children {
		buff += childModel.View()
	}

	return buff
}

func main() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))

	var appSpinner *spinner.Model = &s

	if _, err := tea.NewProgram(model{
		spinner: appSpinner,
		children: []tea.Model{
			NewVpnModel(appSpinner),
			NewHostSelectorModel(appSpinner),
		},
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
