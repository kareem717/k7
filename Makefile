build:
	@go build -tags dev -o bin/k7 main.go  

run: build
	@./bin/k7

install:
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download