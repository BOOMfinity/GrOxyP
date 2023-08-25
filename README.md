# GrOxyP - Go Proxy and VPN Checker

[![Go Reference](https://pkg.go.dev/badge/github.com/BOOMfinity/GrOxyP.svg)](https://pkg.go.dev/github.com/BOOMfinity/GrOxyP)
[![CodeFactor](https://www.codefactor.io/repository/github/boomfinity/groxyp/badge)](https://www.codefactor.io/repository/github/boomfinity/groxyp)

Check if user is behind a VPN or proxy via simple HTTP request.

## Sources

This app is using [X4BNet's list](https://github.com/X4BNet/lists_vpn) of IPs. GrOxyP checks if queried IP is on this
list.

Sources of code are mentioned in the comments.

## Benchmarks

Ran on Windows 11 22631.2262, AMD Ryzen 7 3700X, 32GB RAM 3200MHz, Go 1.21.

[go-wrk](https://github.com/tsliwowicz/go-wrk) benchmark:

- 100 connections, 20 seconds:

```shell
$ go-wrk -c 100 -d 20 "http://localhost:5656/ip?q=194.35.232.123&token=token"

Running 20s test @ http://localhost:5656/ip?q=194.35.232.123&token=token
  100 goroutine(s) running concurrently
1356936 requests in 18.422209262s, 197.99MB read
Requests/sec:           73657.62
Transfer/sec:           10.75MB
Avg Req Time:           1.357632ms
Fastest Request:        0s
Slowest Request:        18.6332ms
Number of Errors:       0
```

Stats (Task Manager, eyeballing): ~30% CPU, ~40MB RAM

- 1 connection, 20 seconds:

```shell
$ go-wrk -c 1 -d 20 "http://localhost:5656/ip?q=194.35.232.123&token=token"

Running 20s test @ http://localhost:5656/ip?q=194.35.232.123&token=token
  1 goroutine(s) running concurrently
243087 requests in 19.0139044s, 35.47MB read
Requests/sec:           12784.70
Transfer/sec:           1.87MB
Avg Req Time:           78.218µs
Fastest Request:        0s
Slowest Request:        4.6746ms
Number of Errors:       0
```

Stats (Task Manager, eyeballing): ~11% CPU, ~38MB RAM

# Installation and usage

## As a server

1. Clone: `git clone https://github.com/BOOMfinity/GrOxyP`.
2. Go to directory: `cd GrOxyP/cmd/groxyp`.
3. Build: `go build`.
4. Set environmental variables as in example:

```sh
  GROXYP_DB_URL = "https://raw.githubusercontent.com/X4BNet/lists_vpn/main/output/datacenter/ipv4.txt"
  GROXYP_DB_UPDATE_INTERVAL = "4h0m0s"
  GROXYP_PORT = 5656,
  GROXYP_TOKEN = "such_a_token_wow"
  GROXYP_DEBUG = false
```

> [!NOTE]
> Port and token are only needed, when you want to spin up an HTTP server.

> [!IMPORTANT]
> Always refer to [X4BNet's repo](https://github.com/X4BNet/lists_vpn) for more information about IP lists. You might
> want to replace `datacenter` with `vpn` in the example above.

5. Run!

HTTP server will be waiting for requests at default port 5656. Query `ip` endpoint like so:

```shell
$ curl http://localhost:5656/ip?q=194.35.232.123&token=such_a_token_wow
{"ip":"194.35.232.123","proxy":true,"rule":"194.35.232.0/22"}
```

Invalid token will cause `401 Unauthorized` messages. Other endpoints will respond with `404` message.

## Programmatically

Use example code below:

```go
package main

import (
	"fmt"
	groxyp "github.com/BOOMfinity/GrOxyP/pkg/client"
	"net"
)

var ipChecker, _ = groxyp.NewClient(groxyp.Config{
	DatabaseDownloadURL:    "https://raw.githubusercontent.com/X4BNet/lists_vpn/main/output/datacenter/ipv4.txt",
	DatabaseUpdateInterval: "8h0m0s",
})

func main() {
	found, reason := ipChecker.FindIP(net.ParseIP("8.8.8.8"))
	fmt.Printf("IP found in the list: %v. IP block: %v", found, reason)
}
```

# Discord support server

Because why not?

[![Discord Widget](https://discordapp.com/api/guilds/1036320104486547466/widget.png?style=banner4)](https://labs.boomfinity.xyz)
