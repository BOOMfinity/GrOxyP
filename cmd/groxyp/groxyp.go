package main

import (
	"fmt"
	"github.com/BOOMfinity/GrOxyP/pkg/config"
	"github.com/BOOMfinity/GrOxyP/pkg/database"
	"github.com/BOOMfinity/GrOxyP/pkg/webserver"
)

func main() {
	// Check envs
	conf, err := config.Get()
	if err != nil {
		panic(err)
	}
	// Downloading fresh database immediately
	err = database.Update(&conf)
	if err != nil {
		panic(err)
		return
	}
	// Updating database "in background" at given interval
	go func() {
		// Starting interval
		err = database.SetUpdateInterval(&conf)
		if err != nil {
			fmt.Println(err)
		}
	}()
	// Starting webserver to listen HTTP queries
	err = webserver.Listen(conf.WebserverPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
