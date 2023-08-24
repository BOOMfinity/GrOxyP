package config

import (
	"fmt"
	"os"
)

// Get reads environmental variables
func Get() (Config, error) {
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
		DatabaseDownloadURL:    os.Getenv("GROXYP_DB_URL"),
		DatabaseUpdateInterval: os.Getenv("GROXYP_DB_UPDATE_INTERVAL"),
		WebserverPort:          os.Getenv("GROXYP_PORT"),
		Debug:                  os.Getenv("GROXYP_DEBUG") == "true",
		Token:                  os.Getenv("GROXYP_TOKEN")}, nil
}

// Config is a structure of environmental variables
type Config struct {
	DatabaseDownloadURL    string
	DatabaseUpdateInterval string
	WebserverPort          string
	Token                  string
	Debug                  bool
}
