package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/ip2location/ip2proxy-go"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "OK\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}
func ip(db *ip2proxy.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		//IP for testing: "103.121.38.138" - should be proxy
		ip := req.FormValue("q")
		proxy, err := db.IsProxy(ip)

		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := apiResponseIpType{
			IP:   ip,
			Type: proxyTypes[(proxy)],
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func Listen(db *ip2proxy.DB, port uint16) error {
	//Source: https://gobyexample.com/http-servers
	http.HandleFunc("/", hello)
	http.HandleFunc("/ip", ip(db))

	fmt.Println(fmt.Sprintf("INFO: Listening on port %v", port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		return err
	}
	return nil
}
