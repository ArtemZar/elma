.PHONY: build lint

default: lint

build:
	go build

lint:
	golangci-lint run -v ./...
