package main

import (
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/config"
	"github.com/BOOMfinity-Developers/GrOxyP/database"
	"github.com/BOOMfinity-Developers/GrOxyP/webserver"
)

func main() {
	err := database.UpdateDatabase(false)
	if err != nil {
		return
	}
	var cfg = config.GetConfig()
	if err != nil {
		return
	}
	err = webserver.Listen(cfg.WebserverPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
