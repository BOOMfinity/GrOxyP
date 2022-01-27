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
)

var cfg = config.GetConfig()
var ranger = cidranger.NewPCTrieRanger()

func UpdateDatabase(disableUpdate bool) error { //arg for debug
	if disableUpdate {
		return nil
	}
	//Source: https://golang.cafe/blog/golang-unzip-file-example.html
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

func SearchIPInDatabase(query string) (bool, string) {
	if containingNetworks, err := ranger.ContainingNetworks(net.ParseIP(query)); len(containingNetworks) > 0 && err == nil {
		network := containingNetworks[0].Network()
		return true, network.String()
	}

	return false, ""
}
