package webserver

// ipEndpointResponse is a structure of /ip endpoint response
type ipEndpointResponse = struct {
	IP    string `json:"ip"`
	Proxy bool   `json:"proxy"`
	Rule  string `json:"rule"`
}
