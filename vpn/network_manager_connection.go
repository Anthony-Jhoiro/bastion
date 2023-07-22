package vpn

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type NetworkManagerConnexion struct {
	ConnectionName string
}

func (c NetworkManagerConnexion) Name() string {
	return c.ConnectionName
}

var ConnectionError = errors.New("fail To connect to VPN")

func (c NetworkManagerConnexion) IsConnectedToVPN() bool {
	cmd := exec.Command("/usr/bin/nmcli", "c", "show", "--active")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return false
	}

	return strings.Contains(out.String(), c.ConnectionName)
}

func (c NetworkManagerConnexion) ConnectToVPN() error {
	cmd := exec.Command("/usr/bin/nmcli", "c", "up", c.ConnectionName)
	return cmd.Run()
}

func (c NetworkManagerConnexion) EnsureConnectedToVpn() error {
	if !c.IsConnectedToVPN() {
		err := c.ConnectToVPN()
		if err != nil {
			return err
		}

		if !c.IsConnectedToVPN() {
			return ConnectionError
		}
	}
	return nil
}
