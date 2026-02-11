.PHONY: build test clean deps fmt lint

BINARY_NAME=nuvadb
BUILD_DIR=bin
CMD_DIR=cmd

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

build:
	@echo "Building NuvaDB..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-server $(CMD_DIR)/server/main.go
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-client $(CMD_DIR)/client/main.go
	@echo "Build complete!"

deps:
	$(GOMOD) download
	$(GOMOD) tidy

test:
	$(GOTEST) -v -race ./...

fmt:
	gofmt -w .

clean:
	rm -rf $(BUILD_DIR)
	$(GOCMD) clean
