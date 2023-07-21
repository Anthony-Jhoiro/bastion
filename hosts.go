package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const cacheFileName = "/home/anthony/.bastion-data.json"

type Host struct {
	Name string
	Ip   string
	Up   bool
}

func (h Host) Title() string       { return h.Name }
func (h Host) Description() string { return h.Ip }
func (h Host) FilterValue() string { return h.Name }

var NmapLogRegex = regexp.MustCompile("Nmap scan report for (.*) \\((.*)\\)")

func ListHostsInNetwork(ipRange string) ([]Host, error) {
	cmd := exec.Command("/usr/bin/nmap", "-sn", ipRange)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	logs := strings.Split(out.String(), "\n")

	hosts := make([]Host, 0)

	for _, logLine := range logs {
		matches := NmapLogRegex.FindStringSubmatch(logLine)
		if matches != nil {
			hosts = append(hosts, Host{
				Name: matches[1],
				Ip:   matches[2],
				Up:   true,
			})
		}

	}

	WriteHostsInCache(hosts)

	return hosts, nil
}

func WriteHostsInCache(hosts []Host) {
	jsonBytes, err := json.Marshal(hosts)
	if err != nil {
		return
	}

	err = os.WriteFile(cacheFileName, jsonBytes, 0664)
	if err != nil {
		return
	}
}

func GetHostsFromCache() []Host {
	fileContent, err := os.ReadFile(cacheFileName)
	if err != nil {
		return nil
	}
	var hosts []Host

	err = json.Unmarshal(fileContent, &hosts)
	if err != nil {
		return nil
	}
	return hosts
}
