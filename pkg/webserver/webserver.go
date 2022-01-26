package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/database"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "OK\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ip(w http.ResponseWriter, req *http.Request) {
	//IP for testing: uk2345.nordvpn.com [194.35.232.123] - should be proxy
	ip := req.FormValue("q")
	proxy, rule := database.SearchIPInDatabase(ip)
	w.Header().Set("Content-Type", "application/json")
	response := apiResponseIpType{
		IP:    ip,
		Proxy: proxy,
		Rule:  rule,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Listen(port uint16) error {
	//Source: https://gobyexample.com/http-servers
	http.HandleFunc("/", hello)
	http.HandleFunc("/ip", ip)

	fmt.Println(fmt.Sprintf("INFO: Listening on port %v", port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		return err
	}
	return nil
}
