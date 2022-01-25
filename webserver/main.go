package webserver

import (
	"GrOxyP/config"
	"encoding/json"
	"fmt"
	"github.com/ip2location/ip2proxy-go"
	"net/http"
)

var cfg = config.GetConfig()

var db, dberr = ip2proxy.OpenDB(fmt.Sprintf("./db/%v", cfg.DatabaseFilename))

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func ip(w http.ResponseWriter, req *http.Request) {
	//IP for testing: "103.121.38.138" - should be proxy
	ip := req.FormValue("q")
	fmt.Println(ip)
	proxy, err := db.IsProxy(ip)

	if err != nil {
		fmt.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := apiResponseIpType{
		IP:   ip,
		Type: proxyTypes[(proxy)],
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func Listen() error {
	//Source: https://gobyexample.com/http-servers
	http.HandleFunc("/", hello)
	http.HandleFunc("/ip", ip)

	err := http.ListenAndServe(":5656", nil)
	if err != nil {
		return err
	}
	return nil
}
