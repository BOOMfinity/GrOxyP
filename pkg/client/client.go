package client

import (
	"fmt"
	"github.com/yl2chen/cidranger"
	"net"
	"net/http"
	"os"
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
// You can use X4BNetDatacenterIPv4ListURL or X4BNetDatacenterIPv4ListURL, or any other compatible list.
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
}

// NewClient creates a new client with the given config and immediately updates its database.
// It will also run automatic updates if an interval is specified.
func NewClient(conf Config) (client *Client, err error) {
	client = &Client{
		Conf:       conf,
		HTTPClient: &http.Client{},
		Database:   cidranger.NewPCTrieRanger(),
	}
	err = client.Update()
	if err != nil {
		return &Client{}, err
	}
	if client.Conf.DatabaseUpdateInterval.String() != "" {
		go func() {
			err = client.runAutoUpdates()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "GrOxyP auto-update error:", err)
			}
		}()
	}
	return client, nil
}

// FindIP checks if a given IP is on the list. If so, returns 'true' and the network.
func (c *Client) FindIP(query net.IP) (found bool, network net.IPNet) {
	if containingNetworks, err := c.Database.ContainingNetworks(query); len(containingNetworks) > 0 && err == nil {
		return true, containingNetworks[0].Network()
	}
	return false, net.IPNet{}
}
