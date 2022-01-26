package main

import (
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/config"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/database"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/webserver"
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
