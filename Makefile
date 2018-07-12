#!/bin/bash

export PKGS=$(shell go list ./... | grep -v vendor/)

all: build-page format test build run

format:
	@echo "HistoryKV: Formatting all golang file."
	@gofmt -w .

test:
	@echo "HistoryKV: Testing... (Probably none, since I'm lazy to create it)"
	@go test -race ${PKGS}

build-page:
	@echo "HistoryKV: Building HTML Page Admin"
	@echo 'package http' > ./src/http/page-admin.go
	@echo '' >> ./src/http/page-admin.go
	@echo 'func (h *HTTP) GetAdminHTML() string{' >> ./src/http/page-admin.go
	@echo '	return `' >> ./src/http/page-admin.go
	@cat "./html/admin.html" | tr -d '\n' | sed -e 's/\s\+/\ /g' >> ./src/http/page-admin.go
	@echo '' >> ./src/http/page-admin.go
	@echo '`' >> ./src/http/page-admin.go
	@echo '}' >> ./src/http/page-admin.go
	@echo "HistoryKV: Building HTML Page Index"
	@echo 'package http' > ./src/http/page-index.go
	@echo '' >> ./src/http/page-index.go
	@echo 'func (h *HTTP) GetIndexHTML() string{' >> ./src/http/page-index.go
	@echo '	return `' >> ./src/http/page-index.go
	@cat "./html/index.html" | tr -d '\n' | sed -e 's/\s\+/\ /g' >> ./src/http/page-index.go
	@echo '' >> ./src/http/page-index.go
	@echo '`' >> ./src/http/page-index.go
	@echo '}' >> ./src/http/page-index.go


build:
	@echo "HistoryKV: Building into ./historykv"
	@go build -o historykv

run:
	@echo "HistoryKV: Trying to running binary"
	@./historykv -config ./historykv.conf

update:
	@echo "HistoryKV: DEP Update"
	@dep ensure -v

ensure-vendor:
	@echo "HistoryKV: DEP Update Vendor"
	@dep ensure -vendor-only -v

clean:
	@echo "HistoryKV: Cleaning"
	@rm ./historykv

build-release:
	@echo "HistoryKV: Make Release"
	@env GOOS=linux GOARCH=386 go build -o release/historykv-linux-386
	@env GOOS=linux GOARCH=amd64 go build -o release/historykv-linux-amd64
	@env GOOS=darwin GOARCH=386 go build -o release/historykv-darwin-386
	@env GOOS=darwin GOARCH=amd64 go build -o release/historykv-darwin-amd64
	@env GOOS=openbsd GOARCH=386 go build -o release/historykv-openbsd-386
	@env GOOS=openbsd GOARCH=amd64 go build -o release/historykv-openbsd-amd64
