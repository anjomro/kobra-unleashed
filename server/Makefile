# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
# GOTEST = $(GOCMD) test
# GOGET = $(GOCMD) get

# Build parameters
BINARY_NAME = kobra-server
BUILD_DIR = dist
TARGET_OS = linux
TARGET_ARCH = arm

# Default target
all: clean build

# Clean target
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Build target
build:
	@echo "Building for testing..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

# Production build target
build-prod:
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v -ldflags="-s -w"

# Test target
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Get dependencies target
get-deps:
	@echo "Getting dependencies..."
	@$(GOGET) -v ./...

.PHONY: all clean build test get-deps
