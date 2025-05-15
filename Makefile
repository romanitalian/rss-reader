APP_NAME = rssreader
MAIN_PACKAGE = ./cmd/rssreader
GO = go
BUILD_DIR = ./build
BINARY = $(BUILD_DIR)/$(APP_NAME)
CLEAR = clear

# Define operating system
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)

.DEFAULT_GOAL := help

.PHONY: all build clean run test tidy check fmt lint help build-linux build-windows build-macos build-all

##@ Development

all: clean build ## Clean and build application

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BINARY) $(MAIN_PACKAGE)
	@echo "Done! Binary file: $(BINARY)"

fmt: ## Format code
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "Done!"

lint: ## Run linter
	@echo "Checking code..."
	@if [ -x "$(GOPATH)/bin/golangci-lint" ]; then \
		$(GOPATH)/bin/golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin; \
		$(GOPATH)/bin/golangci-lint run; \
	fi

tidy: ## Update dependencies
	@echo "Updating dependencies..."
	@$(GO) mod tidy
	@echo "Done!"

##@ Testing

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v ./...

check: fmt lint ## Check and format code

##@ Running

run: build ## Run the application
	@echo "Running $(APP_NAME)..."
	@$(BINARY)

##@ Building

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)/linux
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/linux/$(APP_NAME) $(MAIN_PACKAGE)
	@echo "Done! Binary file: $(BUILD_DIR)/linux/$(APP_NAME)"

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)/windows
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/windows/$(APP_NAME).exe $(MAIN_PACKAGE)
	@echo "Done! Binary file: $(BUILD_DIR)/windows/$(APP_NAME).exe"

build-macos: ## Build for macOS
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)/macos
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/macos/$(APP_NAME) $(MAIN_PACKAGE)
	@echo "Done! Binary file: $(BUILD_DIR)/macos/$(APP_NAME)"

build-all: build-linux build-windows build-macos ## Build for all platforms

##@ Cleaning

clean: ## Clean the build directory
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Done!"

##@ Help

help: ## Display available commands
	@$(CLEAR)
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[0;33m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo "" 