# Makefile for PDFminion

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=pdfminion
VERSION=0.3.2

# Directories
OUTPUT_DIR=build
DIST_DIR=dist

# Build information
BUILDTIME=$(shell date -u +'%Y %b %d %H:%M')

BUILD_OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
BUILD_ARCH=$(shell uname -m | sed 's/x86_64/amd64/' | sed 's/arch64/arm64/')
BUILD_PLATFORM=$(BUILD_OS)-$(BUILD_ARCH)


# Read version from main.go or VERSION file
VERSION=$(shell grep -E "const appVersion = \"[^\"]+\"" cmd/pdfminion/main.go | cut -d'"' -f2)

# Build flags
LDFLAGS=-ldflags "-s -w \
    -X 'pdfminion/internal/domain.buildTime=$(BUILDTIME)' \
    -X 'pdfminion/internal/domain.buildPlatform=$(BUILD_PLATFORM)'"

# Platform specific settings
WINDOWS_AMD64=windows-amd64
LINUX_AMD64=linux-amd64
DARWIN_AMD64=darwin-amd64
DARWIN_ARM64=darwin-arm64

# Install directory
INSTALL_DIR=/usr/local/bin

.PHONY: help default all clean test compile-all run install uninstall release mac \
        compile-windows-amd64 compile-linux-amd64 compile-darwin-amd64 compile-darwin-arm64 \
        package-windows-amd64 package-linux-amd64 package-darwin-amd64 package-darwin-arm64 \
        docker-test-linux lint count install-gotestsum cover

# Help target - displays all available targets with descriptions
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "PDFminion Makefile - Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}' | \
		sort
	@echo ""
	@echo "Build information:"
	@echo "  Version:  $(VERSION)"
	@echo "  Platform: $(BUILD_PLATFORM)"

# Directory targets
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# Directories must exist
# General rule to create directories dynamically
$(OUTPUT_DIR)/%:
	mkdir -p $@

$(DIST_DIR)/%:
	mkdir -p $@

# Default target for current platform only
default: mac ## Build for current platform (alias for 'mac')

# Build for all platforms
all: test compile-all ## Run tests and build for all platforms

# Compile for all platforms
compile-all: ## Compile binaries for all platforms (Windows, Linux, macOS Intel & ARM)
	@echo "Compiling for all platforms..."
	$(MAKE) compile-windows-amd64
	$(MAKE) compile-linux-amd64
	$(MAKE) compile-darwin-amd64
	$(MAKE) compile-darwin-arm64
	@echo "All platforms compiled successfully"

# Build for current platform
compile-windows-amd64: $(OUTPUT_DIR)/$(BINARY_NAME)-$(WINDOWS_AMD64) ## Build for Windows (amd64)
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) \
	-o $(OUTPUT_DIR)/$(BINARY_NAME)-$(WINDOWS_AMD64)/$(BINARY_NAME).exe \
	-v ./cmd/pdfminion

compile-linux-amd64: $(OUTPUT_DIR)/$(BINARY_NAME)-$(LINUX_AMD64) ## Build for Linux (amd64)
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) \
	-o $(OUTPUT_DIR)/$(BINARY_NAME)-$(LINUX_AMD64)/$(BINARY_NAME) \
	-v ./cmd/pdfminion

compile-darwin-arm64: $(OUTPUT_DIR)/$(BINARY_NAME)-$(DARWIN_ARM64) ## Build for macOS Apple Silicon (arm64)
	@echo "Building for macOS (Apple Silicon, arm64)..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) \
	-o $(OUTPUT_DIR)/$(BINARY_NAME)-$(DARWIN_ARM64)/$(BINARY_NAME) \
	-v ./cmd/pdfminion

compile-darwin-amd64: $(OUTPUT_DIR)/$(BINARY_NAME)-$(DARWIN_AMD64) ## Build for macOS Intel (amd64)
	@echo "Building for macOS (Intel, amd64)..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) \
	-o $(OUTPUT_DIR)/$(BINARY_NAME)-$(DARWIN_AMD64)/$(BINARY_NAME) \
	-v ./cmd/pdfminion

lint: ## Run golangci-lint for code quality checks
	@echo "Running golangci-lint..."
	golangci-lint run

# Shortcut for Apple Silicon build
mac: compile-darwin-arm64 ## Build for macOS Apple Silicon (default target)
	@cp $(OUTPUT_DIR)/$(BINARY_NAME)-$(DARWIN_ARM64)/$(BINARY_NAME) ./$(BINARY_NAME)
	@echo "Apple Silicon build available in both $(OUTPUT_DIR)/$(BINARY_NAME)-$(DARWIN_ARM64)/$(BINARY_NAME)"
	@echo "and ./$(BINARY_NAME)"


# Platform specific packaging
package-windows-amd64: compile-windows-amd64 $(DIST_DIR) ## Package Windows build as ZIP
	@echo "Packaging Windows (amd64) build..."
	cd $(OUTPUT_DIR) && \
	zip -r ../$(DIST_DIR)/$(BINARY_NAME)-$(WINDOWS_AMD64)-$(VERSION).zip \
   	$(BINARY_NAME)-$(WINDOWS_AMD64)
	@echo "Windows package created in $(DIST_DIR)"

package-linux-amd64: compile-linux-amd64 $(DIST_DIR) ## Package Linux build as tar.gz
	@echo "Packaging Linux (amd64) build..."
	cd $(OUTPUT_DIR) && \
	tar czf ../$(DIST_DIR)/$(BINARY_NAME)-$(LINUX_AMD64)-$(VERSION).tar.gz \
   	$(BINARY_NAME)-$(LINUX_AMD64)
	@echo "Linux package created in $(DIST_DIR)"

package-darwin-amd64: compile-darwin-amd64 $(DIST_DIR) ## Package macOS Intel build as tar.gz
	@echo "Packaging macOS (amd64) build..."
	cd $(OUTPUT_DIR) && \
   tar czf ../$(DIST_DIR)/$(BINARY_NAME)-$(DARWIN_AMD64)-$(VERSION).tar.gz \
	$(BINARY_NAME)-$(DARWIN_AMD64)
	@echo "macOS package created in $(DIST_DIR)"

package-darwin-arm64: compile-darwin-arm64 $(DIST_DIR) ## Package macOS ARM build as tar.gz
	@echo "Packaging macOS (Apple Silicon) build..."
	cd $(OUTPUT_DIR) && \
	tar czf ../$(DIST_DIR)/$(BINARY_NAME)-$(DARWIN_ARM64)-$(VERSION).tar.gz \
		$(BINARY_NAME)-$(DARWIN_ARM64)
	@echo "macOS (Apple Silicon) package created in $(DIST_DIR)"

# Build all platforms
release: package-windows-amd64 package-linux-amd64 package-darwin-amd64 package-darwin-arm64 ## Build and package for all platforms (release ready)
	@echo "All platform builds completed!"
	@ls -l $(DIST_DIR)


# Add new targets (near the test target):
install-gotestsum: ## Install gotestsum test runner
	go install gotest.tools/gotestsum@latest

test: install-gotestsum ## Run all tests with gotestsum
	gotestsum --format=testdox ./...

count: ## Count lines of code (requires cloc)
	cloc --exclude-dir=vendor,dist,build .

# Target to calculate coverage for all packages
cover: ## Calculate and display test coverage
	@echo "Calculating test coverage for all packages..."
	@go test -coverprofile=coverage.out ./...
	@echo "Package coverage summary:"
	@go tool cover -func=coverage.out | tail -n1 | sed 's/\s\+/ /g' | sed 's/%%/%/'

# Test binaries in containers
docker-test-linux: compile-linux-amd64 ## Test Linux build in Docker container
	docker build -f Dockerfile.linux-test -t pdfminion-linux-test .
	@echo "Testing Linux build..."
	docker run --rm pdfminion-linux-test ./pdfminion --version
	docker run --rm pdfminion-linux-test ./pdfminion --help

clean: ## Clean build artifacts and binaries
	$(GOCLEAN)
	rm -rf $(OUTPUT_DIR)
	rm -rf $(DIST_DIR)
	rm -f ./$(BINARY_NAME)  # Remove the binary from the root directory

run: compile ## Build and run pdfminion
	./$(OUTPUT_DIR)/$(BINARY_NAME)

# Install (Unix-like systems only)
install: compile ## Install pdfminion to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@if [ ! -d "$(INSTALL_DIR)" ]; then \
		echo "Error: Install directory $(INSTALL_DIR) does not exist. Aborting."; \
		exit 1; \
	fi
	@if [ ! -w "$(INSTALL_DIR)" ]; then \
		echo "Error: No write permission for $(INSTALL_DIR). Aborting."; \
		exit 1; \
	fi
	@sudo mv $(OUTPUT_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME) || { \
		echo "Error: Failed to move binary to $(INSTALL_DIR). Aborting."; \
		exit 1; \
	}
	@echo "Installation complete. You can now run '$(BINARY_NAME)' from anywhere."

# Uninstall
uninstall: ## Uninstall pdfminion from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@if [ ! -f "$(INSTALL_DIR)/$(BINARY_NAME)" ]; then \
		echo "Error: $(BINARY_NAME) not found in $(INSTALL_DIR). Aborting."; \
		exit 1; \
	fi
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME) || { \
		echo "Error: Failed to remove $(BINARY_NAME) from $(INSTALL_DIR). Aborting."; \
		exit 1; \
	}
	@echo "Uninstallation complete."
