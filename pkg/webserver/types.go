package webserver

// apiResponseIP is a structure of /ip endpoint response
type apiResponseIP = struct {
	IP    string `json:"ip"`
	Proxy bool   `json:"proxy"`
	Rule  string `json:"rule"`
}
