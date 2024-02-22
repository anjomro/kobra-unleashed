# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
# GOTEST = $(GOCMD) test
# GOGET = $(GOCMD) get

# Build parameters
BINARY_NAME = kobra-server
BUILD_DIR = dist
TARGET_ARCH = armv7

# Default target
all: clean build

# Clean target
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Build target
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

# Test target
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Get dependencies target
get-deps:
	@echo "Getting dependencies..."
	@$(GOGET) -v ./...

.PHONY: all clean build test get-deps
