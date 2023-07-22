package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ullaakut/nmap/v3"
	"log"
	"os"
	"time"
)

const cacheFileName = "/home/anthony/.bastion-data.json"

type Host struct {
	Name string
	Ip   string
	Up   bool
}

func (h Host) Title() string {
	if h.Up {
		return fmt.Sprintf("%v ✅", h.Name)
	}
	return h.Name
}
func (h Host) Description() string { return h.Ip }
func (h Host) FilterValue() string { return h.Name }

func ListHostsInNetwork(ipRanges []string) ([]Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(ipRanges...),
		nmap.WithPorts("22"),
	)
	if err != nil {
		return nil, err
	}

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	hosts := make([]Host, 0, len(result.Hosts))

	for _, nmapHost := range result.Hosts {
		if len(nmapHost.Addresses) == 0 {
			continue
		}

		host := Host{
			Ip: nmapHost.Addresses[0].Addr,
			Up: true,
		}

		if len(nmapHost.Hostnames) == 0 {
			host.Name = "Unknown"
		} else {
			host.Name = nmapHost.Hostnames[0].Name
		}

		hosts = append(hosts, host)
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

	for i, host := range hosts {
		host.Up = false
		hosts[i] = host
	}

	return hosts
}
