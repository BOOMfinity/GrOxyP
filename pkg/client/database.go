package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/yl2chen/cidranger"
	"log"
	"net"
	"strings"
	"time"
)

// Update replaces the database only after the entire downloaded list is valid.
func (c *Client) Update() (err error) {
	if c.Conf.Debug {
		log.Println("INFO: Downloading database...")
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
		log.Println("INFO: Database downloaded. Parsing...")
	}
	database := cidranger.NewPCTrieRanger()

	// Importing
	scanner := bufio.NewScanner(response.Body)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])
		if line == "" {
			continue
		}
		_, network, err := net.ParseCIDR(line)
		if err != nil {
			return fmt.Errorf("invalid CIDR on line %d: %w", lineNumber, err)
		}
		if err := database.Insert(cidranger.NewBasicRangerEntry(*network)); err != nil {
			return fmt.Errorf("inserting CIDR on line %d: %w", lineNumber, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if database.Len() == 0 {
		return errors.New("downloaded database is empty")
	}
	c.mu.Lock()
	c.Database = database
	c.mu.Unlock()
	if c.Conf.Debug {
		log.Printf("INFO: Database parsed. %v entries.\n", database.Len())
	}
	return nil
}

// runAutoUpdates is a simple function to run Update at the set interval
func (c *Client) runAutoUpdates() {
	ticker := time.NewTicker(c.Conf.DatabaseUpdateInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			if c.Conf.Debug {
				log.Println("INFO: Database update started...")
			}
			if err := c.Update(); err != nil {
				log.Printf("GrOxyP auto-update error: %v; keeping the previous database and retrying on the next tick", err)
			}
		}
	}
}
