package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"time"
)

type HostSelectorModel struct {
	tea.Model
	spinner       *spinner.Model
	error         error
	hostsFetched  bool
	hostsFetching bool
	hosts         list.Model
	selectedHost  Host

	width  int
	height int
}

func NewHostSelectorModel(s *spinner.Model) HostSelectorModel {
	marginX, _ := docStyle.GetFrameSize()
	return HostSelectorModel{
		spinner: s,
		hosts:   list.New([]list.Item{}, list.NewDefaultDelegate(), 30-marginX, 16),
	}
}

type SshConnectionFinishedMsg struct {
	err error
}

func startSshConnection(host Host) tea.Cmd {
	c := exec.Command("/usr/bin/ssh", host.Ip)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return SshConnectionFinishedMsg{err}
	})
}

func (m HostSelectorModel) Init() tea.Cmd {
	return listHostsInNetworkFromCache
}

func WithDelay(duration time.Duration, msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(duration)
		return msg
	}
}

type StartDiscoveryMsg struct {
}

func (m HostSelectorModel) onConnectionEstablished(msg ConnectionEstablished) (tea.Model, tea.Cmd) {
	m.hostsFetching = true
	return m, listHostsInNetwork
}

func (m HostSelectorModel) onHostSelected() (tea.Model, tea.Cmd) {
	if m.hostsFetched {
		i, ok := m.hosts.SelectedItem().(Host)
		if ok {
			m.selectedHost = i
		}
		return m, startSshConnection(i)
	}
	return m, nil
}

func (m HostSelectorModel) onListHostsResponse(msg ListHostsResponse) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.error = msg.err
		return m, nil
	}
	marginX, _ := docStyle.GetFrameSize()
	m.hosts = list.New(msg.hosts, list.NewDefaultDelegate(), m.width-marginX, 16)
	m.hostsFetched = true

	if msg.source == AutoDiscovery {
		m.hostsFetching = false
		return m, WithDelay(5*time.Second, StartDiscoveryMsg{})
	}

	return m, nil
}

func (m HostSelectorModel) onSshConnectionFinishedMsg(msg SshConnectionFinishedMsg) (tea.Model, tea.Cmd) {
	m.error = msg.err
	return m, nil
}

func (m HostSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.onHostSelected()
		}
	case ListHostsResponse:
		return m.onListHostsResponse(msg)

	case ConnectionEstablished:
		return m.onConnectionEstablished(msg)

	case StartDiscoveryMsg:
		m.hostsFetching = true
		return m, listHostsInNetwork

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		marginX, marginY := docStyle.GetFrameSize()
		m.hosts.SetSize(m.width-marginX, 16-marginY)

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

func (m HostSelectorModel) View() string {
	buff := ""

	if m.error != nil {
		buff += errorKeyword(fmt.Sprintf("    [ERROR]: %v\n\n", m.error.Error()))
	}

	if m.hostsFetching {
		buff += fmt.Sprintf("  %vLoading hosts in your network \n\n", m.spinner.View())
	}

	return buff + m.hosts.View()
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
