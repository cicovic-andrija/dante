# Makefile for Go module github.com/cicovic-andrija/dante
#

BIN = dantesrv
BIN_DIR = bin

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN) -v ./cmd

.PHONY: clean
clean:
	@go clean
	@rm -rf $(BIN_DIR)

.PHONY: deploy
deploy:
	@./deploy/dev-local.sh
