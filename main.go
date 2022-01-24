package main

import (
	"fmt"
	"github.com/ip2location/ip2proxy-go"
)

var proxyTypes = map[int]string{
	-1: "error",
	0:  "not_proxy",
	1:  "proxy",
	2:  "hosting",
}

func main() {
	/*err := database.UpdateDatabase()
	if err != nil {
		return
	}*/
	db, err := ip2proxy.OpenDB("./db/IP2PROXY-LITE-PX2.BIN")

	if err != nil {
		panic(err)
	}

	ip := "103.121.38.138" //random proxy IP to test with
	proxy, err := db.IsProxy(ip)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Proxy: %v\n", proxyTypes[int(proxy)])

	db.Close()
}
