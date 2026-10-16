.DEFAULT_GOAL := build

BINARY := gobi
BIN_DIR := bin
MAIN_PKG := ./cmd/server
CLIENT := ./tools/client/client.go

HTTP_DIR := tools/http
HTTP_ENV := development
IJHTTP := ijhttp --env-file $(HTTP_DIR)/http-client.env.json -e $(HTTP_ENV)

.PHONY: build clean lint fmt http-get http-patch http-post

# Requires golangci-lint v2 (https://golangci-lint.run/welcome/install/).
lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

build: lint
	go build -o $(BIN_DIR)/$(BINARY) $(MAIN_PKG)

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
