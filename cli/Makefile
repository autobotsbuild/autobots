.PHONY: build test lint clean

BINARY_NAME=autobots
BUILD_DIR=bin
VERSION ?= v0.0.1

build:
	go build -ldflags "-X 'github.com/SheaHawkins/AutoBots/internal/cmd.Version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME) ./

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
