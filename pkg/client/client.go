package client

import (
	"github.com/yl2chen/cidranger"
	"net"
	"net/http"
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

// NewClient creates new client with given config and immediately updates its database.
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
	return client, nil
}

// FindIP checks if given IP is on the list. If so, returns "true" and reason.
func (c *Client) FindIP(query net.IP) (bool, string) {
	if containingNetworks, err := c.Database.ContainingNetworks(query); len(containingNetworks) > 0 && err == nil {
		network := containingNetworks[0].Network()
		return true, network.String()
	}
	return false, ""
}
