package database

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BOOMfinity/GrOxyP/pkg/config"
	"github.com/yl2chen/cidranger"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Defining CIDR checker to check, if given IP is included in given CIDR
var ranger = cidranger.NewPCTrieRanger()

// Update is for downloading database from GitHub to ips.txt and then storing it in memory
func Update(conf *config.Config) error {
	// Source: https://golang.cafe/blog/golang-unzip-file-example.html
	fmt.Println("INFO: Downloading database...")
	response, err := http.Get(conf.DatabaseDownloadURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("received code %v while downloading database", response.StatusCode))
	}
	file, err := os.Create(conf.DatabaseFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	fmt.Println("INFO: Downloading done")
	err = convert(*conf)
	if err != nil {
		return err
	}
	return nil
}

// SetUpdateInterval is simple function to run Update at given interval
func SetUpdateInterval(conf *config.Config) error {
	interval, err := time.ParseDuration(conf.DatabaseUpdateInterval)
	if err != nil {
		log.Fatal(err)
	}
	for range time.Tick(interval) {
		fmt.Println("INFO: Database update started...")
		err := Update(conf)
		if err != nil {
			return err
		}
		fmt.Println("INFO: Database update done")
	}
	return nil
}

// convert converts (wow) downloaded ips.txt file to networks in memory
func convert(conf config.Config) error {
	file, err := os.Open(conf.DatabaseFilename)
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

// FindIP checks if given IP is on the list. If so, returns "true" and reason.
func FindIP(query string) (bool, string) {
	if containingNetworks, err := ranger.ContainingNetworks(net.ParseIP(query)); len(containingNetworks) > 0 && err == nil {
		network := containingNetworks[0].Network()
		return true, network.String()
	}
	return false, ""
}
