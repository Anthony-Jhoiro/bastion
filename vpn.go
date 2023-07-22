package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type VpnModel struct {
	connectionStatus ConnectionStatus
	spinner          *spinner.Model
}

func NewVpnModel(s *spinner.Model) VpnModel {
	return VpnModel{
		connectionStatus: Connecting,
		spinner:          s,
	}
}

type ConnectionEstablished struct{}

func (m VpnModel) Init() tea.Cmd {
	return ensureConnectedToVpn
}

func (m VpnModel) onConnectionStatus(msg ConnectionStatus) (tea.Model, tea.Cmd) {
	m.connectionStatus = msg
	if msg == ConnectionFailed {
		return m, tea.Quit
	}

	if m.connectionStatus == Connected {
		return m, func() tea.Msg {
			return ConnectionEstablished{}
		}
	}
	return m, nil
}

func (m VpnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConnectionStatus:
		return m.onConnectionStatus(msg)
	}
	return m, nil
}

func (m VpnModel) View() string {
	buff := ""

	loadingPrefix := fmt.Sprintf("  %v", m.spinner.View())
	defaultPrefix := "    "

	if m.connectionStatus == Connecting {
		return buff + fmt.Sprintf("%vConnecting to %v\n\n", loadingPrefix, infoKeyword(VpnName))
	}
	if m.connectionStatus == ConnectionFailed {
		return buff + fmt.Sprintf("%vFail to connect to %v\n\n", defaultPrefix, errorKeyword(VpnName))
	}

	return buff + fmt.Sprintf("%vConnection to %v established\n\n", defaultPrefix, successKeyword(VpnName))
}

func ensureConnectedToVpn() tea.Msg {
	if err := EnsureConnectedToVpn(); err != nil {
		return ConnectionFailed
	}
	return Connected
}
