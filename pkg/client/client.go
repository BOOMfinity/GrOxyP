package client

import (
	"fmt"
	"github.com/yl2chen/cidranger"
	"net"
	"net/http"
	"sync"
	"time"
)

// X4BNetVPNIPv4ListURL is the shorthand for the primary remote list of VPNs.
// The URL was confirmed to be online at the release time; however, it may break in the future.
// As per X4BNet's repo:
//
// "This list is strictly just known VPN networks.
//
//	Small overlap with datacenter networks is possible (e.g., if it isn't possible to separate) however,
//	most datacenters will not be in this list"
const X4BNetVPNIPv4ListURL = "https://raw.githubusercontent.com/X4BNet/lists_vpn/main/output/vpn/ipv4.txt"

// X4BNetDatacenterIPv4ListURL is the shorthand for the primary remote list of datacenters *AND* VPNs.
// The URL was confirmed to be online at the release time; however, it may break in the future.
// As per X4BNet's repo:
//
// "This list is for VPNs and Datacenters.
//
//	Anything that is 'not an eyeball network' directly."
//
// Cool definition btw.: https://en.wikipedia.org/wiki/Eyeball_network
const X4BNetDatacenterIPv4ListURL = "https://raw.githubusercontent.com/X4BNet/lists_vpn/main/output/datacenter/ipv4.txt"

// Config is a structure of environmental variables
//
// DatabaseDownloadURL: URL to the remote list.
// You can use X4BNetVPNIPv4ListURL or X4BNetDatacenterIPv4ListURL, or any other compatible list.
// DatabaseUpdateInterval: update interval; no duration disables updates
// WebserverPort: port when in standalone mode
// WebserverToken: auth token when in standalone mode
// Debug: debug flag, useful to... debug!
type Config struct {
	DatabaseDownloadURL    string
	DatabaseUpdateInterval time.Duration
	WebserverPort          string
	WebserverToken         string
	Debug                  bool
}

type Client struct {
	Conf       Config
	Database   cidranger.Ranger
	HTTPClient *http.Client
	mu         sync.RWMutex
	stop       chan struct{}
	stopOnce   sync.Once
	updates    sync.WaitGroup
}

// NewClient creates a new client with the given config and immediately updates its database.
// It will also run automatic updates if an interval is specified.
func NewClient(conf Config) (client *Client, err error) {
	if conf.DatabaseUpdateInterval < 0 {
		return nil, fmt.Errorf("database update interval must not be negative")
	}
	client = &Client{
		Conf:       conf,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Database:   cidranger.NewPCTrieRanger(),
		stop:       make(chan struct{}),
	}
	err = client.Update()
	if err != nil {
		return nil, err
	}
	if client.Conf.DatabaseUpdateInterval > 0 {
		client.updates.Add(1)
		go func() {
			defer client.updates.Done()
			client.runAutoUpdates()
		}()
	}
	return client, nil
}

// Close stops automatic updates and waits for any ongoing download to finish.
func (c *Client) Close() {
	if c.stop == nil {
		return
	}
	c.stopOnce.Do(func() { close(c.stop) })
	c.updates.Wait()
}

// FindIP checks if a given IP is on the list. If so, returns 'true' and the network.
func (c *Client) FindIP(query net.IP) (found bool, network net.IPNet) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if query == nil || c.Database == nil {
		return false, net.IPNet{}
	}
	if containingNetworks, err := c.Database.ContainingNetworks(query); len(containingNetworks) > 0 && err == nil {
		return true, containingNetworks[0].Network()
	}
	return false, net.IPNet{}
}
