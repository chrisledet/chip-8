DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES = $(shell go list ./...)

all: clean deps build run

deps:
	@echo "-> Installing dependencies"
	@go get -d -v ./... $(DEPS)

test:
	@echo "-> Testing..."
	@go test ./...

clean:
	@echo "-> Cleaning..."
	@rm -f ./chip-8

build:
	@echo "-> Building..."
	@go build

run:
	@echo "-> Running..."
	@./chip-8

.PHONY: all clean
