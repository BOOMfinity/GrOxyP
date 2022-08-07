package config

// Config is a structure of config.js file
type Config struct {
	DatabaseFilename       string `json:"databaseFilename"`
	DatabaseDownloadURL    string `json:"databaseDownloadURL"`
	DatabaseUpdateInterval string `json:"databaseUpdateInterval"`
	WebserverPort          string `json:"webserverPort"`
	Token                  string `json:"token"`
}
