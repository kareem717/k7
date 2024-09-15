# Simple Makefile for a Go project

build:
	@go build -o bin/kareem main.go

run: build
	@./bin/proj

install:
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download

test:
	@echo "Testing..."
	@go test ./... -v

watch:
	@air -c .air.toml

