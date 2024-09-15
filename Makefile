build:
	@go build -tags dev -o bin/github.com/kareem717/k7 main.go  

run: build
	@./bin/github.com/kareem717/k7

install:
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download