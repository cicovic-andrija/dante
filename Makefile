# Makefile for Go module github.com/cicovic-andrija/dante
#

BIN = dantesrv
BIN_DIR = bin
BIN_USER = $(HOME)/bin

.PHONY: install
install: build
	@cp $(BIN_DIR)/$(BIN) $(BIN_USER)/$(BIN)

.PHONY: build
build: prebuild
	go build -o $(BIN_DIR)/$(BIN) -v ./cmd

.PHONY: prebuild
prebuild:
	@mkdir -p $(BIN_DIR)

.PHONY: clean
clean:
	@go clean
	@rm -rf $(BIN_DIR)

.PHONY: tidy
tidy:
	@go mod tidy
