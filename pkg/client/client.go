package client

import (
	"fmt"
	"github.com/yl2chen/cidranger"
	"net"
	"net/http"
	"os"
)

// Config is a structure of environmental variables
type Config struct {
	DatabaseDownloadURL    string
	DatabaseUpdateInterval string
	WebserverPort          string
	WebserverToken         string
	Debug                  bool
}

type Client struct {
	Conf       Config
	Database   cidranger.Ranger
	HTTPClient *http.Client
}

// NewClient creates new client with the given config and immediately updates its database. It will also run automatic updates, if interval is specified.
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
	if client.Conf.DatabaseUpdateInterval != "" {
		go func() {
			err = client.runAutoUpdates()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "GrOxyP auto-update error:", err)
			}
		}()
	}
	return client, nil
}

// FindIP checks if a given IP is on the list. If so, returns `true` and the reason (IP block).
func (c *Client) FindIP(query net.IP) (bool, string) {
	if containingNetworks, err := c.Database.ContainingNetworks(query); len(containingNetworks) > 0 && err == nil {
		network := containingNetworks[0].Network()
		return true, network.String()
	}
	return false, ""
}
