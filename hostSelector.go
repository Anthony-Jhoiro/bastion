package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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

func (m HostSelectorModel) Init() tea.Cmd {
	return listHostsInNetwork
}

func WithDelay(duration time.Duration, msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(duration)
		return msg
	}
}

type StartDiscoveryMsg struct {
}

func (m HostSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			if m.hostsFetched {
				i, ok := m.hosts.SelectedItem().(Host)
				if ok {
					m.selectedHost = i
				}
				return m, startSshConnection(i)
			}
		}
	case ListHostsResponse:
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

	case StartDiscoveryMsg:
		m.hostsFetching = true
		return m, listHostsInNetwork

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		marginX, marginY := docStyle.GetFrameSize()
		m.hosts.SetSize(m.width-marginX, 16-marginY)

	case sshConnectionFinishedMsg:
		if msg.err != nil {
			m.error = msg.err
			return m, tea.Quit
		}
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

	if m.hostsFetching || !m.hostsFetched {
		buff += fmt.Sprintf("  %vLoading hosts in your network \n\n", m.spinner.View())
	}

	return buff + m.hosts.View()
}
