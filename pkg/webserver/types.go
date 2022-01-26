package webserver

type apiResponseIpType = struct {
	IP    string `json:"ip"`
	Proxy bool   `json:"proxy"`
	Rule  string `json:"rule"`
}
