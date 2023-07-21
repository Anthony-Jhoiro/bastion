package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

const VpnName = "hydra"

var VpnConnectionError = errors.New("fail To connect to VPN")

func IsConnectedToVPN() bool {
	cmd := exec.Command("/usr/bin/nmcli", "c", "show", "--active")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return false
	}

	return strings.Contains(out.String(), VpnName)
}

func ConnectToVPN() error {
	cmd := exec.Command("/usr/bin/nmcli", "c", "up", VpnName)
	return cmd.Run()
}

func EnsureConnectedToVpn() error {
	if !IsConnectedToVPN() {
		err := ConnectToVPN()
		if err != nil {
			return err
		}

		if !IsConnectedToVPN() {
			return VpnConnectionError
		}
	}
	return nil
}
