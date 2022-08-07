package webserver

import (
	"fmt"
	"github.com/BOOMfinity-Developers/GrOxyP/pkg/database"
	"github.com/segmentio/encoding/json"
	"net/http"
	"os"
)

// notfound returns "OK" on every non-existing endpoint
func notfound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(404)
}

// ip returns queried IP, if queried IP is behind a proxy or VPN and which network has been blocked (reason/rule)
func ip(w http.ResponseWriter, req *http.Request) {
	// IP for testing: uk2345.nordvpn.com [194.35.232.123] - should be proxy
	ip := req.FormValue("q")
	token := req.FormValue("token")
	if token != os.Getenv("GROXYP_TOKEN") {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := fmt.Fprintf(w, "401 Unauthorized")
		if err != nil {
			return
		}
	} else {
		proxy, rule := database.FindIP(ip)
		w.Header().Set("Content-Type", "application/json")
		response := ipEndpointResponse{
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
}

// Listen starts HTTP server for IP queries. Available endpoints: /ip. Usage is in README.
func Listen(port string) error {
	//Source: https://gobyexample.com/http-servers
	http.HandleFunc("/", notfound)
	http.HandleFunc("/ip", ip)

	fmt.Println(fmt.Sprintf("INFO: Listening on port %v", port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		return err
	}
	return nil
}
