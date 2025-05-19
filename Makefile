.PHONY: build install clean test run

APP_NAME = shai

# Build directories
CMD_DIR = cmd/cli
BUILD_DIR = build
BIN_DIR = $(GOPATH)/bin

GO_BUILD_FLAGS = -v

# Determine the operating system
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	INSTALL_PATH = /usr/local/bin
else ifeq ($(UNAME_S),Linux)
	INSTALL_PATH = /usr/local/bin
else
	INSTALL_PATH = $(BIN_DIR)
endif

all: build

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_PATH)/$(APP_NAME)..."
	@mkdir -p $(INSTALL_PATH)
	cp $(BUILD_DIR)/$(APP_NAME) $(INSTALL_PATH)/
	@echo "Installation complete. Make sure $(INSTALL_PATH) is in your PATH."

run: build
	./$(BUILD_DIR)/$(APP_NAME)

test:
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	@echo "Cleaned build directory"

help:
	@echo "Makefile for $(APP_NAME)"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the application"
	@echo "  make install  - Install the application to $(INSTALL_PATH)"
	@echo "  make run      - Build and run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make help     - Show this help message"
