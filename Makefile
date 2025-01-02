
OUTPUT = 403unlocker
MAIN = cmd/403unlockercli/main.go
BIN_DIR = ~/.local/bin

.DEFAULT_GOAL := help

.PHONY: help lint build test clean install uninstall


help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'


lint:
	@golangci-lint run


build:
	@go build -o $(OUTPUT) $(MAIN)


test: 
	@go test ./...


clean: 
	@rm -f $(OUTPUT)


install: build 
	@echo "Installing $(OUTPUT) to $(BIN_DIR)..."
	@install -m 755 $(OUTPUT) $(BIN_DIR)


uninstall:
	@echo "Removing $(OUTPUT) from $(BIN_DIR)..."
	@rm -f $(BIN_DIR)/$(OUTPUT)
