package hosts

import (
	"github.com/Anthony-Jhoiro/bastion/hosts/discovery"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) listHostsInNetwork() tea.Msg {
	hosts, err := m.discovery.DiscoverHosts()
	if err != nil {
		return discovery.ListHostsResponse{
			Err: err,
		}
	}

	lst := make([]list.Item, len(hosts))
	for i, host := range hosts {
		lst[i] = host
	}

	return discovery.ListHostsResponse{
		Hosts:  lst,
		Err:    err,
		Source: discovery.AutoDiscovery,
	}
}

func (m model) listHostsInNetworkFromCache() tea.Msg {
	hosts := m.discovery.GetHostsFromCache()
	if hosts == nil {
		return nil
	}
	lst := make([]list.Item, len(hosts))
	for i, host := range hosts {
		lst[i] = host
	}

	return discovery.ListHostsResponse{
		Hosts:  lst,
		Source: discovery.Cache,
	}
}
