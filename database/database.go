package database

import (
	"GrOxyP/config"
	"GrOxyP/unzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var cfg = config.GetConfig()

func UpdateDatabase(disableUpdate bool) error { //arg for debug
	if disableUpdate {
		return nil
	}
	//Source: https://golang.cafe/blog/golang-unzip-file-example.html
	fmt.Println("INFO: Downloading database...")
	URL := fmt.Sprintf("https://www.ip2location.com/download?token=%v&file=%v", cfg.DatabaseToken, cfg.DatabaseCode)
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("received code %v while downloading database", response.StatusCode))
	}
	file, err := os.Create("db.zip")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	fmt.Println("INFO: Downloading done")

	err = unpackDatabase()
	if err != nil {
		return err
	}

	return nil
}

func unpackDatabase() error {
	fmt.Println("INFO: Unzipping...")
	err := unzip.Run("db.zip", "db")
	if err != nil {
		return err
	}
	fmt.Println("INFO: Unzipping done")
	return nil
}
