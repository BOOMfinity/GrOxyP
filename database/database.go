package database

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/config"
	"io"
	"net"
	"net/http"
	"os"
)

var cfg = config.GetConfig()
var nets []*net.IPNet

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
			nets = append(nets, currNet)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func SearchIPInDatabase(query string) (bool, string) {
	q := net.ParseIP(query)
	for _, currNet := range nets {
		if currNet.Contains(q) {
			return true, currNet.String()
		}
	}
	return false, ""
}
