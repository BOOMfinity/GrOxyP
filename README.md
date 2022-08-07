# GrOxyP - Go Proxy and VPN Checker

[![Go Reference](https://pkg.go.dev/badge/github.com/BOOMfinity/GrOxyP.svg)](https://pkg.go.dev/github.com/BOOMfinity/GrOxyP)
[![CodeFactor](https://www.codefactor.io/repository/github/boomfinity/groxyp/badge)](https://www.codefactor.io/repository/github/boomfinity/groxyp)
[![BCH compliance](https://bettercodehub.com/edge/badge/BOOMfinity/GrOxyP?branch=master)](https://bettercodehub.com/)

Check if user is behind a VPN or proxy via simple HTTP request.

## Sources

This app is using [X4BNet's list](https://github.com/X4BNet/lists_vpn) of IPs. GrOxyP checks if queried IP is on this
list.

Sources of code are mentioned in the comments.

## Benchmarks

Ran on Windows 11, AMD Ryzen 7 3700X, 32GB RAM 3200MHz.

[go-wrk](https://github.com/tsliwowicz/go-wrk) benchmark:

- 100 connections, 20 seconds:

```shell
$ go-wrk -c 100 -d 20 http://localhost:5656/ip?q=194.35.232.123

Running 20s test @ http://localhost:5656/ip?q=194.35.232.123
  100 goroutine(s) running concurrently
574077 requests in 17.079976921s, 83.76MB read
Requests/sec:           33611.11
Transfer/sec:           4.90MB
Avg Req Time:           2.975206ms
Fastest Request:        0s
Slowest Request:        32.4233ms
Number of Errors:       0
# Stats: ~20% CPU, ~50MB RAM
```

- 1 connection, 20 seconds:

```shell
$ go-wrk -c 1 -d 20 http://localhost:5656/ip?q=194.35.232.123

Running 20s test @ http://localhost:5656/ip?q=194.35.232.123
  1 goroutine(s) running concurrently
283966 requests in 18.9446991s, 41.43MB read
Requests/sec:           14989.21
Transfer/sec:           2.19MB
Avg Req Time:           66.714µs
Fastest Request:        0s
Slowest Request:        3.641ms
Number of Errors:       0
# Stats: ~10% CPU, ~38MB RAM
```

# Installation and usage

1. Clone: `git clone https://github.com/BOOMfinity/GrOxyP`.
2. Go to directory: `cd GrOxyP/cmd/groxyp`.
3. Build: `go build`.
4. Set environmental variables as in example:

```shell
  GROXYP_DB_URL = "https://raw.githubusercontent.com/X4BNet/lists_vpn/main/ipv4.txt"
  GROXYP_DB_FILE = "ips.txt"
  GROXYP_DB_UPDATE_INTERVAL = "4h0m0s"
  GROXYP_PORT = 5656,
  GROXYP_TOKEN: = "such_a_token_wow"
```

5. Run!

HTTP server will be waiting for requests at default port 5656. Query `ip` endpoint like so:

```shell
$ curl http://localhost:5656/ip?q=194.35.232.123&token=such_a_token_wow
{"ip":"194.35.232.123","proxy":true,"rule":"194.35.232.0/22"}
```

Invalid token will cause `401 Unauthorized` messages. Other endpoints should respond with `404` message.
