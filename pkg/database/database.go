package database

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/config"
	"github.com/yl2chen/cidranger"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// Getting config
var cfg = config.GetConfig()

// Defining CIDR checker to check, if given IP is included in given CIDR
var ranger = cidranger.NewPCTrieRanger()

// UpdateDatabase is for downloading database from GitHub to ips.txt and then storing it in memory
func UpdateDatabase(disableUpdate bool) error {
	// If disableUpdate is true, application will NOT update its database. Useful for debug or offline mode
	if disableUpdate {
		return nil
	}
	// Source: https://golang.cafe/blog/golang-unzip-file-example.html
	fmt.Println("INFO: Downloading database...")
	response, err := http.Get("https://raw.githubusercontent.com/X4BNet/lists_vpn/main/ipv4.txt")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("received code %v while downloading database", response.StatusCode))
	}
	file, err := os.Create(cfg.DatabaseFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	fmt.Println("INFO: Downloading done")
	err = convertDatabase()
	if err != nil {
		return err
	}
	return nil
}

// SetUpdateInterval is simple function to run UpdateDatabase at given interval
func SetUpdateInterval(d time.Duration, disableUpdate bool) error {
	for range time.Tick(d) {
		fmt.Println("INFO: Database update started...")
		err := UpdateDatabase(disableUpdate)
		if err != nil {
			return err
		}
		fmt.Println("INFO: Database update done")
	}
	return nil
}

// convertDatabase converts downloaded ips.txt file to networks in memory
func convertDatabase() error {
	file, err := os.Open(cfg.DatabaseFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if _, currNet, err := net.ParseCIDR(scanner.Text()); err == nil {
			err := ranger.Insert(cidranger.NewBasicRangerEntry(*currNet))
			if err != nil {
				fmt.Printf("Error while inserting CIDR to database: %v\n", err.Error())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// SearchIPInDatabase checks if given IP is on the list. If so, returns "true" and reason.
func SearchIPInDatabase(query string) (bool, string) {
	if containingNetworks, err := ranger.ContainingNetworks(net.ParseIP(query)); len(containingNetworks) > 0 && err == nil {
		network := containingNetworks[0].Network()
		return true, network.String()
	}
	return false, ""
}
