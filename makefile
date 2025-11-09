.PHONY: build clean install test run help

BINARY_NAME=netly
VERSION=1.0.0
BUILD_DIR=build

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo '${GREEN}Available targets:${NC}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${YELLOW}%-15s${NC} %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "${GREEN}Building ${BINARY_NAME}...${NC}"
	@go build -o ${BINARY_NAME} -ldflags "-s -w" .
	@echo "${GREEN}Build complete: ${BINARY_NAME}${NC}"

build-all: ## Build for all platforms
	@echo "${GREEN}Building for all platforms...${NC}"
	@mkdir -p ${BUILD_DIR}
	@GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 -ldflags "-s -w" .
	@GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 -ldflags "-s -w" .
	@GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe -ldflags "-s -w" .
	@GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 -ldflags "-s -w" .
	@GOOS=darwin GOARCH=arm64 go build -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 -ldflags "-s -w" .
	@echo "${GREEN}All builds complete in ${BUILD_DIR}/${NC}"

clean: ## Clean build files
	@echo "${YELLOW}Cleaning...${NC}"
	@rm -rf ${BINARY_NAME} ${BUILD_DIR}
	@echo "${GREEN}Clean complete${NC}"

install: build ## Install the binary to /usr/local/bin
	@echo "${GREEN}Installing ${BINARY_NAME}...${NC}"
	@sudo mv ${BINARY_NAME} /usr/local/bin/
	@echo "${GREEN}Installation complete${NC}"

uninstall: ## Uninstall the binary
	@echo "${YELLOW}Uninstalling ${BINARY_NAME}...${NC}"
	@sudo rm -f /usr/local/bin/${BINARY_NAME}
	@echo "${GREEN}Uninstall complete${NC}"

run: build ## Build and run the binary
	@./${BINARY_NAME}

deps: ## Download dependencies
	@echo "${GREEN}Downloading dependencies...${NC}"
	@go mod download
	@go mod tidy
	@echo "${GREEN}Dependencies updated${NC}"

test: ## Run tests
	@echo "${GREEN}Running tests...${NC}"
	@go test -v ./...

fmt: ## Format code
	@echo "${GREEN}Formatting code...${NC}"
	@go fmt ./...
	@echo "${GREEN}Format complete${NC}"

lint: ## Run linter
	@echo "${GREEN}Running linter...${NC}"
	@golangci-lint run || echo "${YELLOW}Note: Install golangci-lint to use this target${NC}"

dev: ## Run in development mode with auto-reload
	@echo "${GREEN}Starting development mode...${NC}"
	@air || go run .
