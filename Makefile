.DEFAULT_GOAL := build

BINARY := gobi
BIN_DIR := bin
MAIN_PKG := ./cmd/server
CLIENT := ./tools/client/client.go

HTTP_DIR := tools/http
HTTP_ENV := development
IJHTTP := ijhttp --env-file $(HTTP_DIR)/http-client.env.json -e $(HTTP_ENV)

CERT_HOST ?= gobi.local
CERT_DIR := tools/certs

.PHONY: build clean lint fmt certs http-get http-patch http-post

# Requires golangci-lint v2 (https://golangci-lint.run/welcome/install/).
lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

build: lint
	go build -o $(BIN_DIR)/$(BINARY) $(MAIN_PKG)

# Requires mkcert (https://github.com/FiloSottile/mkcert). Override the host with
# `make certs CERT_HOST=other.local`. Output goes to tools/certs/<host>.pem and <host>-key.pem.
certs:
	mkcert -install
	mkdir -p $(CERT_DIR)
	mkcert -cert-file $(CERT_DIR)/$(CERT_HOST).pem -key-file $(CERT_DIR)/$(CERT_HOST)-key.pem $(CERT_HOST)

clean:
	rm -rf $(BIN_DIR)

client:
	go run $(CLIENT)

# The http-* targets send requests to a running server (start it with ./bin/gobi).
http-get:
	$(IJHTTP) $(HTTP_DIR)/get.http

http-patch:
	$(IJHTTP) $(HTTP_DIR)/patch.http

http-post:
	$(IJHTTP) $(HTTP_DIR)/post.http
