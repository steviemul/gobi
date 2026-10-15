# gobi

A small HTTP/1.1 web server written in Go, built from the TCP socket up as a learning project.

## Description

gobi doesn't use Go's `net/http` package. It accepts raw TCP connections with `net`, parses requests and builds responses itself, and manages connection lifetime directly. The goal is to understand how an HTTP server works under the hood, and to grow it into a simple general-purpose server that can serve static files, proxy to backends and host multiple virtual hosts.

### What works today

- HTTP/1.0 and HTTP/1.1 request parsing, with proper `400` and `413` responses for malformed or oversized requests
- Keep-alive connections (default for 1.1, opt-in for 1.0), with pipelined requests
- Idle, read and write timeouts, so slow or abandoned clients don't hold connections open
- A built-in status page, and different resource types in responses
- Optional YAML configuration, falling back to defaults when no file is given or the file doesn't exist

### Planned

The config file already describes routes, but the handlers that act on them are not implemented yet:

- **Static files**: a route with `root` serves files from that directory
- **Reverse proxy**: a route with `proxy` forwards requests to another server
- **Virtual hosts**: a route with `host` matches on the request's `Host` header

## Getting started

```sh
make build                              # builds bin/gobi
./bin/gobi                              # run with default config (port 8080)
./bin/gobi -config tools/config/default.yaml
```

`make http-get`, `make http-post` and `make http-patch` send sample requests to a running server using JetBrains' `ijhttp` client. The requests live in `tools/http`.

## Configuration

The only command-line flag is `-config`, the path to an optional YAML file. If the flag is omitted or the file doesn't exist, defaults are used. A file that exists but can't be parsed is an error.

```yaml
port: 8090
status: true
idleTimeoutSeconds: 60   # max wait between requests on a keep-alive connection
readTimeoutSeconds: 10   # max time to receive one request
writeTimeoutSeconds: 10  # max time to send one response
routes:                  # evaluated in the order listed
  - path: /www
    root: tools/www
  - path: /backend
    proxy: http://backend.local
  - host: gobi.local
    root: tools/www/gobi
```

| Key | Default | Meaning |
|---|---|---|
| `port` | `8080` | TCP port to listen on |
| `idleTimeoutSeconds` | `60` | Keep-alive idle timeout |
| `readTimeoutSeconds` | `10` | Time allowed to receive a request |
| `writeTimeoutSeconds` | `10` | Time allowed to send a response |
| `routes` | none | Match rules (`host`, `path`) and an action (`root` or `proxy`) |

See [`tools/config/default.yaml`](tools/config/default.yaml) for a full example.

## Project layout

```
cmd/server/            entry point and flag parsing
internal/config/       YAML config loading and defaults
internal/server/       TCP listener and connection handling
internal/server/http/  request parsing and response building
internal/server/handler/ request handlers
tools/                 sample config, HTTP requests, embedded static content, test client
```
