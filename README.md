# gobi

[![CI](https://github.com/steviemul/gobi/actions/workflows/ci.yml/badge.svg)](https://github.com/steviemul/gobi/actions/workflows/ci.yml)

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

## Tools

Besides Go, you need:

| Tool | Needed for | Install |
|---|---|---|
| [`golangci-lint`](https://golangci-lint.run) v2 | `make build`, `make demo`, `make lint` and `make fmt` (linting and formatting, configured in `.golangci.yml`) | `brew install golangci-lint`, or see the [install docs](https://golangci-lint.run/welcome/install/) |
| [`mkcert`](https://github.com/FiloSottile/mkcert) | `make certs` and `make demo` (generates one locally-trusted TLS certificate for `localhost` and `gobi.local`, written to `tools/certs`) | `brew install mkcert`, or see the [install docs](https://github.com/FiloSottile/mkcert#installation) |
| [`ijhttp`](https://www.jetbrains.com/help/idea/http-client-cli.html) | `make http-get`, `make http-post` and `make http-patch` (optional, only to send the sample requests) | See the [JetBrains HTTP Client CLI docs](https://www.jetbrains.com/help/idea/http-client-cli.html) |

## Getting started

```sh
make build                              # builds bin/gobi
./bin/gobi                              # run with default config (port 8080)
./bin/gobi -config tools/config/default.yaml
```

`./bin/gobi -serve` is a shortcut for browsing the current directory: it serves a listing of it at <http://localhost:8080/>, with no TLS and no status page. It can't be combined with `-config`.

`make demo` builds the server and runs it with `tools/config/demo.yaml`, which serves a small multi-page site from `tools/www/demo`. It runs `make certs` first if the certificates don't exist, then open <https://localhost:8443/demo/>.

`make http-get`, `make http-post` and `make http-patch` send sample requests to a running server using JetBrains' `ijhttp` client. The requests live in `tools/http`.

## Configuration

The only command-line flag is `-config`, the path to an optional YAML file. If the flag is omitted or the file doesn't exist, defaults are used (a missing file also logs a warning). A file that exists but can't be parsed is an error.

```yaml
port: 8090
status: true
idleTimeoutSeconds: 60   # max wait between requests on a keep-alive connection
readTimeoutSeconds: 10   # max time to receive one request
writeTimeoutSeconds: 10  # max time to send one response
tls:                     # optional, see TLS below
  enabled: true
  port: 8443
  cert: tools/certs/gobi.pem
  key: tools/certs/gobi-key.pem
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
| `tls` | disabled | HTTPS listener settings, see [TLS](#tls) |
| `routes` | none | Match rules (`host`, `path`) and an action (`root` or `proxy`) |

See [`tools/config/default.yaml`](tools/config/default.yaml) for a full example.

### TLS

TLS is off by default. When enabled, gobi listens for HTTPS on its own port in addition to the plain HTTP port.

**Before enabling it, run `make certs`.** This uses [`mkcert`](https://github.com/FiloSottile/mkcert) to create a locally-trusted certificate for `localhost` and `gobi.local`, and writes `tools/certs/gobi.pem` and `tools/certs/gobi-key.pem`. If the files named in `cert` and `key` can't be read, the HTTPS listener is not started and an error is logged. The certs are git-ignored, so each checkout needs to generate its own.

```yaml
tls:
  enabled: true
  port: 8443
  cert: tools/certs/gobi.pem
  key: tools/certs/gobi-key.pem
```

| Key | Default | Meaning |
|---|---|---|
| `enabled` | `false` | Start the HTTPS listener |
| `port` | `8443` | TCP port for HTTPS |
| `cert` | none | Path to the PEM certificate |
| `key` | none | Path to the PEM private key |
| `force` | `false` | When `true`, the plain HTTP listener answers every request with a `302` redirect to the same path on HTTPS. The server refuses to start if the TLS listener can't start |

The redirect keeps the request's host, path and query, and swaps in the HTTPS port. A request with no `Host` header can't be redirected and gets a `400`.

Once the server is running, try `curl https://localhost:8443/` or `https://gobi.local:8443/` (the latter needs both `127.0.0.1 gobi.local` and `::1 gobi.local` in `/etc/hosts`; with only the IPv4 line, macOS waits about 5 seconds for a multicast DNS answer for the IPv6 lookup).

## Roadmap

Ideas for what to build next are in [`docs/roadmap.md`](docs/roadmap.md).

## Project layout

```
cmd/server/            entry point and flag parsing
internal/config/       YAML config loading and defaults
internal/server/       TCP listener and connection handling
internal/server/http/  request parsing and response building
internal/server/handler/ request handlers
docs/                  roadmap and other notes
tools/                 sample config, HTTP requests, embedded static content, test client
```
