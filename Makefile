.PHONY: build clean install test run-example help

BINARY_NAME=gomcc
INSTALL_PATH=/usr/local/bin

help:
	@echo "Available targets:"
	@echo "  build        - Build the gomcc binary"
	@echo "  clean        - Remove built binaries"
	@echo "  install      - Install gomcc to $(INSTALL_PATH) (requires sudo)"
	@echo "  test         - Run tests"
	@echo "  run-example  - Run with example configuration in dry-run mode"
	@echo "  help         - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	@echo "Clean complete"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete"

test:
	@echo "Running tests..."
	go test -v ./...

run-example: build
	@echo "Running example configuration in dry-run mode..."
	./$(BINARY_NAME) -f example-config.json -r -v
