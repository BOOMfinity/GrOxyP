package webserver

var proxyTypes = map[int8]string{
	-1: "error",
	0:  "not_proxy",
	1:  "proxy",
	2:  "hosting",
}

type apiResponseIpType = struct {
	IP   string `json:"ip"`
	Type string `json:"type"`
}
