package vpn

type Connection interface {
	IsConnectedToVPN() bool
	ConnectToVPN() error
	EnsureConnectedToVpn() error

	Name() string
}
