# Roadmap: possible next features

gobi is a learning project: the aim is something usable, not production ready. These are
ideas for what to build next, ordered roughly by how much Go they teach against how useful
they are. Notes about the current code were accurate when this was written, so check them
before starting.

## Suggested order

1. Methods and `HEAD`, plus the `Date` header
2. Access logging
3. Graceful shutdown
4. Reverse proxy
5. Conditional requests

The first two are quick wins that make the rest easier to debug.

## Best learning value

### Reverse proxy

The main planned feature. `internal/server/handler/proxyHandler.go` is still a stub, and a
route with `proxy` is already accepted in the config.

- **Go practice:** `net.Dial`, `io.Copy`, goroutines copying bytes in both directions,
  `context` for timeouts and cancellation.
- **HTTP practice:** forwarding `Host`, adding `X-Forwarded-For`, stripping hop-by-hop headers
  (`Connection`, `Keep-Alive`, `Transfer-Encoding`, and so on), mapping backend failures to
  `502` and `504`.

### Graceful shutdown

`cmd/server/main.go` has no signal handling, so Ctrl+C kills in-flight requests.

- Use `os/signal` and `context` to catch `SIGINT` and `SIGTERM`.
- Track active connections with a `sync.WaitGroup`. On shutdown, stop accepting, then wait
  for connections to drain, with a deadline.
- `GobiListener.Shutdown` already closes the listeners, so this builds on it.

### Access logging

One line per request: method, path, status, bytes sent, duration and client address.

- A good fit for a small wrapper around the dispatcher call in `handleConnection`.
- Would have made several things we hit visible straight away, such as redirects, 404s and
  stray `/favicon.ico` requests.

### Tests

Only `redirect_test.go` and `fileHandler_test.go` exist today.

- **Parser:** table tests feeding `http.ParseRequest` a `bufio.NewReader(strings.NewReader(...))`.
- **Connections:** an integration test on a random port (`:0`) covering keep-alive, pipelining,
  timeouts and the `400` and `413` responses. `net.Pipe()` is an alternative for tests that
  don't need a real socket.
- **File handler:** path traversal attempts, missing files, index files, content types.

## Correctness gaps

Small changes, each a real HTTP-compliance point.

- **Methods.** The file handler doesn't look at the method, so `POST /demo/` likely serves
  the index. Return `405 Method Not Allowed` with an `Allow` header, and handle `HEAD` by
  sending the headers (including `Content-Length`) without the body.
- **`Date` and `Server` headers.** HTTP/1.1 origin servers are expected to send `Date`.
  Headers are also written from a Go map in `GobiResponse.Bytes`, so their order changes
  between runs. An ordered header type would fix that.
- **Chunked request bodies.** `parseBody` appears to handle only `Content-Length`. Supporting
  `Transfer-Encoding: chunked` is a good way to learn the format.
- **Subdirectory index files.** Only the route root maps to `index.html`, so `/demo/docs/`
  returns `404`. The `/demo` to `/demo/` redirect (`directoryRedirect` in `fileHandler.go`)
  could be extended to any directory on disk: redirect to the slash form, then serve its
  `index.html`.

## Features that make it more usable

- **Conditional requests.** `Last-Modified` or `ETag` with `If-Modified-Since` and
  `If-None-Match`, answering `304 Not Modified`. The biggest practical win for a static server.
  Uses `os.Stat` and the `time` package.
- **Range requests.** `Range` and `206 Partial Content`, which video seeking needs. A good
  exercise in `io.SectionReader`. Also advertise `Accept-Ranges: bytes`.
- **gzip.** Compress text responses when the client sends `Accept-Encoding: gzip`, using
  `compress/gzip`. Add `Vary: Accept-Encoding`.
- **Config validation and reload.** Fail fast on bad config instead of logging and skipping
  (a route with both `root` and `proxy` currently only logs). A `SIGHUP` reload would teach
  safe config swapping with `sync/atomic` or a mutex.
- **Connection limits.** Cap concurrent connections with a buffered-channel semaphore, a
  classic Go pattern that protects the server under load.
- **`tls.force`.** Already done. See the README for how it behaves.

## Probably not worth it

HTTP/2, WebSockets and CGI. They are large, and the learning value falls off compared to
the items above.
