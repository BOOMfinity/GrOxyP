package main

import (
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/config"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/database"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/webserver"
	"log"
	"time"
)

func main() {
	// Getting config from config.json
	var cfg = config.Get()
	// Downloading fresh database immediately
	err := database.Update(false)
	if err != nil {
		return
	}
	// Updating database "in background" at given interval
	go func() {
		// Parsing duration
		interval, err := time.ParseDuration(cfg.DatabaseUpdateInterval)
		if err != nil {
			log.Fatal(err)
		}
		// Starting interval
		err = database.SetUpdateInterval(interval, false)
		if err != nil {
			fmt.Println(err)
		}
	}()
	// Starting webserver to listen HTTP queries
	err = webserver.Listen(cfg.WebserverPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
