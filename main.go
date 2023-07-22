package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"os"
	"os/exec"
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
	quitting         bool
	connectionStatus ConnectionStatus
	spinner          spinner.Model
	error            error

	hostSelectorModel HostSelectorModel
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		ensureConnectedToVpn,
		listHostsInNetworkFromCache,
	)
}

type sshConnectionFinishedMsg struct {
	err error
}

func startSshConnection(host Host) tea.Cmd {
	c := exec.Command("/usr/bin/ssh", host.Ip)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return sshConnectionFinishedMsg{err}
	})
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
			m.hostSelectorModel.hostsFetching = true
			return m, listHostsInNetwork
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.hostSelectorModel.spinner = m.spinner

		return m, cmd
	}

	m2, cmd := m.hostSelectorModel.Update(msg)
	m.hostSelectorModel = m2.(HostSelectorModel)

	if cmd != nil {
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
	buff += m.hostSelectorModel.View()

	return buff
}

func main() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))

	if _, err := tea.NewProgram(model{
		connectionStatus:  Connecting,
		spinner:           s,
		hostSelectorModel: NewHostSelectorModel(s),
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

type ListHostsResponseSource int

const (
	Cache         ListHostsResponseSource = iota
	AutoDiscovery                         = iota
)

type ListHostsResponse struct {
	hosts  []list.Item
	source ListHostsResponseSource
	err    error
}

func listHostsInNetwork() tea.Msg {
	hosts, err := ListHostsInNetwork([]string{"10.0.0.0/24"})
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
		hosts:  lst,
		err:    err,
		source: AutoDiscovery,
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
		hosts:  lst,
		source: Cache,
	}
}
