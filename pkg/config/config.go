package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// GetConfig reads config.json file and returns parsed settings
func GetConfig() Config {
	//Source: https://github.com/MrBoombastic/GoProdukcji
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
	}
	return configuration
}
