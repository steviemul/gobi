.DEFAULT_GOAL := build

BINARY := gobi
BIN_DIR := bin
MAIN_PKG := ./cmd/server
CLIENT := ./tools/client/client.go

HTTP_DIR := tools/http
HTTP_ENV := development
IJHTTP := ijhttp --env-file $(HTTP_DIR)/http-client.env.json -e $(HTTP_ENV)

CERT_HOSTS ?= localhost gobi.local
CERT_NAME := gobi
CERT_DIR := tools/certs
DEMO_CONFIG := tools/config/demo.yaml

.PHONY: build clean lint fmt certs demo http-get http-patch http-post

# Requires golangci-lint v2 (https://golangci-lint.run/welcome/install/).
lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

build: lint
	go build -o $(BIN_DIR)/$(BINARY) $(MAIN_PKG)

# Requires mkcert (https://github.com/FiloSottile/mkcert). One certificate covers every host
# in CERT_HOSTS; override with `make certs CERT_HOSTS="localhost other.local"`.
# Output goes to tools/certs/gobi.pem and gobi-key.pem.
certs:
	mkcert -install
	mkdir -p $(CERT_DIR)
	mkcert -cert-file $(CERT_DIR)/$(CERT_NAME).pem -key-file $(CERT_DIR)/$(CERT_NAME)-key.pem $(CERT_HOSTS)

# Builds and runs the server with the demo site; generates the certs first if they're missing.
# Then browse to https://localhost:8443/demo/.
demo: build $(CERT_DIR)/$(CERT_NAME).pem
	./$(BIN_DIR)/$(BINARY) -config $(DEMO_CONFIG)

$(CERT_DIR)/$(CERT_NAME).pem:
	$(MAKE) certs

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
