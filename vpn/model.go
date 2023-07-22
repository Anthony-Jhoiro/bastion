package vpn

import (
	"bastion/colors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectionStatus = int

const (
	Connected        ConnectionStatus = iota
	ConnectionFailed ConnectionStatus = iota
	Connecting       ConnectionStatus = iota
)

type model struct {
	connexion        Connection
	connectionStatus ConnectionStatus
	spinner          *spinner.Model
}

func NewVpnModel(s *spinner.Model, c Connection) tea.Model {
	return model{
		connectionStatus: Connecting,
		spinner:          s,
		connexion:        c,
	}
}

type ConnectionEstablished struct{}

func (m model) Init() tea.Cmd {
	return m.ensureConnectedToVpn
}

func (m model) onConnectionStatus(msg ConnectionStatus) (tea.Model, tea.Cmd) {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConnectionStatus:
		return m.onConnectionStatus(msg)
	}
	return m, nil
}

func (m model) View() string {
	buff := ""

	loadingPrefix := fmt.Sprintf("  %v", m.spinner.View())
	defaultPrefix := "    "

	if m.connectionStatus == Connecting {
		return buff + fmt.Sprintf("%vConnecting to %v\n\n", loadingPrefix, colors.InfoKeyword(m.connexion.Name()))
	}
	if m.connectionStatus == ConnectionFailed {
		return buff + fmt.Sprintf("%vFail to connect to %v\n\n", defaultPrefix, colors.ErrorKeyword(m.connexion.Name()))
	}

	return buff + fmt.Sprintf("%vConnection to %v established\n\n", defaultPrefix, colors.SuccessKeyword(m.connexion.Name()))
}
