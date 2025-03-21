package main

import (
	"fmt"
	"github.com/BOOMfinity/GrOxyP/pkg/client"
	"os"
	"time"
)

// getConfig reads environmental variables
func getConfig() (c client.Config, err error) {
	if os.Getenv("GROXYP_DB_URL") == "" {
		return client.Config{}, fmt.Errorf("GROXYP_DB_URL is not set")
	}
	if os.Getenv("GROXYP_PORT") == "" {
		return client.Config{}, fmt.Errorf("GROXYP_PORT is not set")
	}
	if os.Getenv("GROXYP_TOKEN") == "" {
		return client.Config{}, fmt.Errorf("GROXYP_TOKEN is not set")
	}
	duration := os.Getenv("GROXYP_DB_UPDATE_INTERVAL")
	var durationParsed time.Duration
	if duration != "" {
		durationParsed, err = time.ParseDuration(duration)
		if err != nil {
			return client.Config{}, fmt.Errorf("error while parsing duration: %+v", err)
		}
	}

	return client.Config{
		DatabaseDownloadURL:    os.Getenv("GROXYP_DB_URL"),
		DatabaseUpdateInterval: durationParsed,
		WebserverPort:          os.Getenv("GROXYP_PORT"),
		Debug:                  os.Getenv("GROXYP_DEBUG") == "true",
		WebserverToken:         os.Getenv("GROXYP_TOKEN")}, nil
}

func main() {
	// Check envs
	conf, err := getConfig()
	if err != nil {
		panic(err)
	}
	// Spinning up a new client
	c, err := client.NewClient(conf)
	if err != nil {
		panic(err)
	}
	// Starting webserver to listen to HTTP queries
	err = c.StartServer()
	if err != nil {
		panic(err)
		return
	}
}
