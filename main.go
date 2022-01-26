package main

import (
	"GrOxyP/config"
	"GrOxyP/database"
	"GrOxyP/webserver"
	"fmt"
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
