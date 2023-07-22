package hosts

import (
	"bastion/colors"
	"bastion/hosts/discovery"
	"bastion/vpn"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"time"
)

type model struct {
	tea.Model
	spinner        *spinner.Model
	error          error
	hostsFetched   bool
	hostsFetching  bool
	hosts          list.Model
	selectedHost   discovery.Host
	readyToConnect bool

	discovery discovery.Discovery
}

func NewHostSelectorModel(s *spinner.Model, discoveryStrategy discovery.Discovery) tea.Model {
	return model{
		spinner:   s,
		hosts:     list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 16),
		discovery: discoveryStrategy,
	}
}

type SshConnectionFinishedMsg struct {
	err error
}

func startSshConnection(host discovery.Host) tea.Cmd {
	c := exec.Command("/usr/bin/ssh", host.Ip)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return SshConnectionFinishedMsg{err}
	})
}

func (m model) Init() tea.Cmd {
	return m.listHostsInNetworkFromCache
}

func WithDelay(duration time.Duration, msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(duration)
		return msg
	}
}

type StartDiscoveryMsg struct {
}

func (m model) onConnectionEstablished() (tea.Model, tea.Cmd) {
	m.hostsFetching = true
	m.readyToConnect = true
	return m, m.listHostsInNetwork
}

func (m model) onHostSelected() (tea.Model, tea.Cmd) {
	if m.hostsFetched {
		i, ok := m.hosts.SelectedItem().(discovery.Host)
		if ok {
			m.selectedHost = i
		}
		return m, startSshConnection(i)
	}
	return m, nil
}

func (m model) onListHostsResponse(msg discovery.ListHostsResponse) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.error = msg.Err
		return m, nil
	}
	m.hosts = list.New(msg.Hosts, list.NewDefaultDelegate(), 30, 16)
	m.hostsFetched = true

	if msg.Source == discovery.AutoDiscovery {
		m.hostsFetching = false
		return m, WithDelay(5*time.Second, StartDiscoveryMsg{})
	}

	return m, nil
}

func (m model) onSshConnectionFinishedMsg(msg SshConnectionFinishedMsg) (tea.Model, tea.Cmd) {
	m.error = msg.err
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.onHostSelected()
		}
	case discovery.ListHostsResponse:
		return m.onListHostsResponse(msg)

	case vpn.ConnectionEstablished:
		return m.onConnectionEstablished()

	case StartDiscoveryMsg:
		m.hostsFetching = true
		return m, m.listHostsInNetwork

	case SshConnectionFinishedMsg:
		return m.onSshConnectionFinishedMsg(msg)
	}

	if m.hostsFetched {
		var cmd tea.Cmd
		m.hosts, cmd = m.hosts.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	buff := ""

	if m.error != nil {
		buff += colors.ErrorKeyword(fmt.Sprintf("    [ERROR]: %v\n\n", m.error.Error()))
	}

	if m.hostsFetching {
		buff += fmt.Sprintf("  %vLoading hosts in your network \n\n", m.spinner.View())
	}

	if m.readyToConnect {
		buff += m.hosts.View()
	}
	return buff

}
