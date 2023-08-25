package client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// IpEndpointResponse is a structure of /ip endpoint response
type IpEndpointResponse = struct {
	IP    string `json:"ip"`
	Proxy bool   `json:"proxy"`
	Rule  string `json:"rule"`
}

// StartServer starts HTTP server for IP queries. Available endpoints: /ip. Usage is in README.
func (c *Client) StartServer() error {
	//Source: https://gobyexample.com/http-servers
	http.HandleFunc("/", notfoundEndpoint)
	http.HandleFunc("/ip", ipEndpoint(c))

	fmt.Println(fmt.Sprintf("INFO: Listening on port %v", c.Conf.WebserverPort))
	err := http.ListenAndServe(fmt.Sprintf(":%v", c.Conf.WebserverPort), nil)
	if err != nil {
		return err
	}
	return nil
}

// notfoundEndpoint returns "OK" on every non-existing endpoint
func notfoundEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(404)
}

// ipEndpoint returns queried IP, if queried IP is behind a proxy or VPN and which network has been blocked (reason/rule)
func ipEndpoint(c *Client) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// IP for testing: uk2345.nordvpn.com [194.35.232.123] - should be proxy
		ip := req.FormValue("q")
		token := req.FormValue("token")
		if token != c.Conf.WebserverToken {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := fmt.Fprintf(w, "401 Unauthorized")
			if err != nil {
				return
			}
		} else {
			proxy, rule := c.FindIP(net.ParseIP(ip))
			w.Header().Set("Content-Type", "application/json")
			response := IpEndpointResponse{
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
}
