package main

import (
	"GrOxyP/database"
	"GrOxyP/webserver"
)

func main() {
	err := database.UpdateDatabase(true)
	if err != nil {
		return
	}

	err = webserver.Listen()
	if err != nil {
		return
	}
}
