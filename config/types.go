package config

type Config struct {
	DatabaseCode     string `json:"databaseCode"`
	DatabaseToken    string `json:"databaseToken"`
	DatabaseFilename string `json:"databaseFilename"`
	WebserverPort    uint16 `json:"webserverPort"`
}
