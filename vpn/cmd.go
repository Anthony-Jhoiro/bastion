package vpn

import tea "github.com/charmbracelet/bubbletea"

func (m model) ensureConnectedToVpn() tea.Msg {
	if err := m.connexion.EnsureConnectedToVpn(); err != nil {
		return ConnectionFailed
	}
	return Connected
}
