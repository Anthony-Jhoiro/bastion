package main

import (
	"bastion/banner"
	"bastion/hosts"
	"bastion/hosts/discovery"
	"bastion/hosts/discovery/nmap"
	"bastion/template"
	"bastion/vpn"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

const cacheFileName = "/home/anthony/.bastion-data.json"

func main() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))

	var appSpinner *spinner.Model = &s

	if _, err := tea.NewProgram(template.NewTemplateModel(
		appSpinner,
		banner.NewBannerModel(),
		vpn.NewVpnModel(appSpinner, vpn.NetworkManagerConnexion{ConnectionName: "hydra"}),
		hosts.NewHostSelectorModel(
			appSpinner,
			discovery.Discovery{
				Strategy: nmap.AutoDiscovery{
					Networks: []string{"10.0.0.0/24", "10.0.1.0/24"},
					Ports:    []string{"22"},
				},
				CacheLocation: cacheFileName,
			},
		),
	), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
