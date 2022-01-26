package main

import (
	"GrOxyP/config"
	"GrOxyP/database"
	"GrOxyP/webserver"
	"fmt"
	"github.com/ip2location/ip2proxy-go"
)

func main() {
	err := database.UpdateDatabase(false)
	if err != nil {
		return
	}
	var cfg = config.GetConfig()
	db, err := ip2proxy.OpenDB(fmt.Sprintf("./db/%v", cfg.DatabaseFilename))
	if err != nil {
		return
	}
	err = webserver.Listen(db, cfg.WebserverPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
