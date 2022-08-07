package config

import (
	"os"
)

// Get reads environmental variables and returns them as Config
func Get() Config {
	return Config{
		DatabaseFilename:       os.Getenv("GROXYP_DB_FILE"),
		DatabaseDownloadURL:    os.Getenv("GROXYP_DB_URL"),
		DatabaseUpdateInterval: os.Getenv("GROXYP_DB_UPDATE_INTERVAL"),
		WebserverPort:          os.Getenv("GROXYP_PORT"),
		Token:                  os.Getenv("GROXYP_TOKEN"),
	}
}
