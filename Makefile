.DEFAULT_GOAL := build

BINARY := gobi
BIN_DIR := bin
MAIN_PKG := ./cmd/server

.PHONY: build clean

build:
	go build -o $(BIN_DIR)/$(BINARY) $(MAIN_PKG)

clean:
	rm -rf $(BIN_DIR)
