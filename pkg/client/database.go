package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/yl2chen/cidranger"
	"log"
	"net"
	"time"
)

// Update will flush and repopulate database with freshly downloaded list.
func (c *Client) Update() (err error) {
	if c.Conf.Debug {
		fmt.Println("INFO: Downloading database...")
	}
	response, err := c.HTTPClient.Get(c.Conf.DatabaseDownloadURL)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("received code %v while downloading database", response.StatusCode))
	}
	if c.Conf.Debug {
		fmt.Println("INFO: Database downloaded. Parsing...")
	}
	// Flushing entries
	c.Database = cidranger.NewPCTrieRanger()

	// Importing
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		if _, currNet, err := net.ParseCIDR(scanner.Text()); err == nil {
			err := c.Database.Insert(cidranger.NewBasicRangerEntry(*currNet))
			if err != nil {
				fmt.Printf("Error while inserting CIDR to database: %v\n", err.Error())
			}
		}
	}
	if c.Conf.Debug {
		fmt.Printf("INFO: Database parsed. %v entries.\n", c.Database.Len())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// runAutoUpdates is simple function to run Update at set interval
func (c *Client) runAutoUpdates() error {
	interval, err := time.ParseDuration(c.Conf.DatabaseUpdateInterval)
	if err != nil {
		log.Fatal(err)
	}
	for range time.Tick(interval) {
		if c.Conf.Debug {
			fmt.Println("INFO: Database update started...")
		}
		err := c.Update()
		if err != nil {
			return err
		}
	}
	return nil
}
