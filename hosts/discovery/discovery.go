package discovery

import (
	"context"
	"encoding/json"
	"os"
	"time"
)

type Strategy interface {
	ListHosts(ctx context.Context) ([]Host, error)
}

type Discovery struct {
	Strategy      Strategy
	CacheLocation string
}

func (d Discovery) DiscoverHosts() ([]Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	foundedHosts, err := d.Strategy.ListHosts(ctx)

	if err != nil {
		return nil, err
	}

	d.writeHostsInCache(foundedHosts)

	return foundedHosts, nil
}

func (d Discovery) writeHostsInCache(hosts []Host) {
	jsonBytes, err := json.Marshal(hosts)
	if err != nil {
		return
	}

	err = os.WriteFile(d.CacheLocation, jsonBytes, 0664)
	if err != nil {
		return
	}
}

func (d Discovery) GetHostsFromCache() []Host {
	fileContent, err := os.ReadFile(d.CacheLocation)
	if err != nil {
		return nil
	}
	var hostsFromCache []Host

	err = json.Unmarshal(fileContent, &hostsFromCache)
	if err != nil {
		return nil
	}

	for i, host := range hostsFromCache {
		host.Up = false
		hostsFromCache[i] = host
	}

	return hostsFromCache
}
