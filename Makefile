.PHONY: test

PACKAGES = $(shell go list ./... | grep -v -e . -e mocks | tr '\n' ',')

build:
	@go build -o bin/ src/cmd/main.go

run: build
	.bin/

test:
	@if [ -f coverage.out ]; then rm coverage.out; fi;
	@echo ">> running unit test and calculate coverage"
	@go test ./... -cover -coverprofile=coverage.out -covermode=count -coverpkg=$(PACKAGES)
	@go tool cover -func=coverage.out