package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
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
	help           = termenv.Style{}.Foreground(color("241")).Styled
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
	quitting         bool
	connectionStatus ConnectionStatus
	spinner          spinner.Model
	error            error
	hostsFetched     bool
	hosts            list.Model
	width            int
	height           int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		ensureConnectedToVpn,
		listHostsInNetworkFromCache,
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

	case ConnectionStatus:
		m.connectionStatus = msg
		if msg == ConnectionFailed {
			return m, tea.Quit
		}
		if m.connectionStatus == Connected {
			return m, listHostsInNetwork
		}
		return m, nil

	case ListHostsResponse:
		if msg.err != nil {
			m.error = msg.err
			return m, nil
		}
		h, _ := docStyle.GetFrameSize()
		m.hosts = list.New(msg.hosts, list.NewDefaultDelegate(), m.width-h, 16)
		m.hostsFetched = true

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.hostsFetched {
			h, _ := docStyle.GetFrameSize()
			m.hosts.SetSize(msg.Width-h, 16)
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.hostsFetched {
		var cmd tea.Cmd
		m.hosts, cmd = m.hosts.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	buff := fmt.Sprintf("%v\n\n", hydraBanner)

	loadingPrefix := fmt.Sprintf("  %v", m.spinner.View())
	defaultPrefix := "    "

	if m.connectionStatus == Connecting {
		return buff + fmt.Sprintf("%vConnecting to %v\n", loadingPrefix, infoKeyword(VpnName))
	}
	if m.connectionStatus == ConnectionFailed {
		return buff + fmt.Sprintf("%vFail to connect to %v\n", defaultPrefix, errorKeyword(VpnName))
	}

	buff += fmt.Sprintf("%vConnection to %v established\n\n", defaultPrefix, successKeyword(VpnName))

	if m.error != nil {
		return buff + errorKeyword(fmt.Sprintf("%vError: %v\n\n", defaultPrefix, m.error))
	}

	if !m.hostsFetched {
		return buff + fmt.Sprintf("%vSearching for hosts\n", loadingPrefix)
	} else {
		return buff + m.hosts.View()
	}
}

func main() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))

	if _, err := tea.NewProgram(model{
		connectionStatus: Connecting,
		spinner:          s,
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func ensureConnectedToVpn() tea.Msg {
	if err := EnsureConnectedToVpn(); err != nil {
		return ConnectionFailed
	}
	return Connected
}

type ListHostsResponse struct {
	hosts []list.Item
	err   error
}

func listHostsInNetwork() tea.Msg {
	hosts, err := ListHostsInNetwork("10.0.0.0/24")
	if err != nil {
		return ListHostsResponse{
			err: err,
		}
	}

	lst := make([]list.Item, len(hosts))
	for i, host := range hosts {
		lst[i] = host
	}

	return ListHostsResponse{
		hosts: lst,
		err:   err,
	}
}

func listHostsInNetworkFromCache() tea.Msg {
	hosts := GetHostsFromCache()
	if hosts == nil {
		return nil
	}
	lst := make([]list.Item, len(hosts))
	for i, host := range hosts {
		lst[i] = host
	}

	return ListHostsResponse{
		hosts: lst,
	}
}
