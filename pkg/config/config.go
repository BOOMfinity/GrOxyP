package config

import (
	"fmt"
	"os"
)

// Get reads environmental variables
func Get() (Config, error) {
	if os.Getenv("GROXYP_DB_FILE") == "" {
		return Config{}, fmt.Errorf("GROXYP_DB_FILE is not set")
	}
	if os.Getenv("GROXYP_DB_URL") == "" {
		return Config{}, fmt.Errorf("GROXYP_DB_URL is not set")
	}
	if os.Getenv("GROXYP_DB_UPDATE_INTERVAL") == "" {
		return Config{}, fmt.Errorf("GROXYP_DB_UPDATE_INTERVAL is not set")
	}
	if os.Getenv("GROXYP_PORT") == "" {
		return Config{}, fmt.Errorf("GROXYP_PORT is not set")
	}
	if os.Getenv("GROXYP_TOKEN") == "" {
		return Config{}, fmt.Errorf("GROXYP_TOKEN is not set")
	}
	return Config{
		DatabaseFilename:       os.Getenv("GROXYP_DB_FILE"),
		DatabaseDownloadURL:    os.Getenv("GROXYP_DB_URL"),
		DatabaseUpdateInterval: os.Getenv("GROXYP_DB_UPDATE_INTERVAL"),
		WebserverPort:          os.Getenv("GROXYP_PORT"),
		Token:                  os.Getenv("GROXYP_TOKEN")}, nil
}

// Config is a structure of environmental variables
type Config struct {
	DatabaseFilename       string
	DatabaseDownloadURL    string
	DatabaseUpdateInterval string
	WebserverPort          string
	Token                  string
}
