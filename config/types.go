package config

type Config struct {
	DatabaseFilename    string `json:"databaseFilename"`
	DatabaseDownloadURL string `json:"databaseDownloadURL"`
	WebserverPort       uint16 `json:"webserverPort"`
}
