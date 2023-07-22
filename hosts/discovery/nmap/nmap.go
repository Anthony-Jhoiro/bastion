package nmap

import (
	"bastion/hosts/discovery"
	"context"
	"github.com/Ullaakut/nmap/v3"
	"log"
)

type AutoDiscovery struct {
	Networks []string
	Ports    []string
}

func (ad AutoDiscovery) ListHosts(ctx context.Context) ([]discovery.Host, error) {

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(ad.Networks...),
		nmap.WithPorts(ad.Ports...),
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

	foundHosts := make([]discovery.Host, 0, len(result.Hosts))

	for _, nmapHost := range result.Hosts {
		if len(nmapHost.Addresses) == 0 {
			continue
		}

		host := discovery.Host{
			Ip: nmapHost.Addresses[0].Addr,
			Up: true,
		}

		if len(nmapHost.Hostnames) == 0 {
			host.Name = "Unknown"
		} else {
			host.Name = nmapHost.Hostnames[0].Name
		}

		foundHosts = append(foundHosts, host)
	}
	return foundHosts, nil
}
