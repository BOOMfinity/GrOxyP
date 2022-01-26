# GrOxyP - Go Proxy

Check if IP is a VPN or proxy. Written in Go.

## Sources

This app is using [X4BNet's list](https://github.com/X4BNet/lists_vpn) of IPs. GrOxyP checks, if queried IP is on this
list.

Sources of code are mentioned in the comments.

## Benchmarks

Run on Windows 11, AMD Ryzen 7 3700X, 32GB RAM 3200MHz.

- [go-wrk](https://github.com/tsliwowicz/go-wrk) benchmark
    - 100 connections, 20 seconds:

```shell
$ go-wrk -c 100 -d 20 http://localhost:5656/ip?q=194.35.232.123

Running 20s test @ http://localhost:5656/ip?q=194.35.232.123
100 goroutine(s) running concurrently
396026 requests in 17.544609038s, 57.79MB read
Requests/sec:           22572.52
Transfer/sec:           3.29MB
Avg Req Time:           4.430165ms
Fastest Request:        0s
Slowest Request:        32.7456ms
Number of Errors:       0
```

About 50% of usage of CPU and about 20MB of RAM.

- 1 connection, 20 seconds:

```shell
$ go-wrk -c 1 -d 20 http://localhost:5656/ip?q=194.35.232.123

Running 20s test @ http://localhost:5656/ip?q=194.35.232.123
  1 goroutine(s) running concurrently
83425 requests in 19.6422217s, 12.17MB read
Requests/sec:           4247.23
Transfer/sec:           634.60KB
Avg Req Time:           235.447µs
Fastest Request:        0s
Slowest Request:        3.1675ms
Number of Errors:       0
```

About 10% of usage of CPU and about 14 MB of RAM.

# Installation and usage

1. Clone: `git clone https://github.com/BOOMfinity-Developers/GrOxyP`.
2. Go to directory: `cd GrOxyP`.
3. Build: `go build`.
4. Copy `config.json.example` and rename to `config.json` and modify it (if you wish).
5. Run!

HTTP server will be ready for requests at default port 5656. Query endpoint `ip` like so:

```shell
$ curl http://localhost:5656/ip?q=194.35.232.123
{"ip":"194.35.232.123","proxy":true,"rule":"194.35.232.0/22"}
```

Any other endpoint should respond with `OK` message.