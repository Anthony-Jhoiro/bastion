package discovery

import "github.com/charmbracelet/bubbles/list"

type ListHostsResponseSource int

const (
	Cache         ListHostsResponseSource = iota
	AutoDiscovery                         = iota
)

type ListHostsResponse struct {
	Hosts  []list.Item
	Source ListHostsResponseSource
	Err    error
}
