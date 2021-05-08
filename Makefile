# Makefile for Go module github.com/cicovic-andrija/rac
#

BINLOC = bin
BINUSER = $(HOME)/bin
BINF = rac

.PHONY: install
install: build
	cp $(BINLOC)/$(BINF) $(BINUSER)/$(BINF)

.PHONY: build
build: prebuild
	go build -o $(BINLOC)/$(BINF) -v ./cmd

.PHONY: prebuild
prebuild:
	mkdir -p $(BINLOC)

.PHONY: clean
clean:
	go clean
	rm -rf $(BINLOC)

.PHONY: tidy
tidy:
	go mod tidy
