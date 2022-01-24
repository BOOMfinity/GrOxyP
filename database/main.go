package database

import (
	"errors"
	"fmt"
	"github.com/MrBoombastic/GrOxyP/unzip"
	"io"
	"net/http"
	"os"
)

func UpdateDatabase() error {
	//Source: https://golang.cafe/blog/golang-unzip-file-example.html
	fmt.Println("Downloading...")

	databaseCode := "PX2LITEBIN"
	token := "nope"
	URL := fmt.Sprintf("https://www.ip2location.com/download?token=%v&file=%v", token, databaseCode)
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("received code %v while downloading datavbase", response.StatusCode))
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

	fmt.Println("Unzipping...")
	err = unpackDatabase()
	if err != nil {
		return err
	}

	return nil
}

func unpackDatabase() error {
	err := unzip.Run("db.zip", "db")
	if err != nil {
		return err
	}
	return nil
}
