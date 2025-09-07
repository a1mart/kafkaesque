# Project Variables
BINARY_NAME := kafkaesque
SRC_DIR := .
BUILD_DIR := bin
SCRIPTS_DIR := pkg/scripts
PROTO_DIR := pkg/proto/huginn
PROTO_OUT := internal/generated
SWAGGER_OUT := pkg/swagger
THIRD_PARTY_DIR := pkg/third_party
GO_FILES := $(shell find $(SRC_DIR) -type f -name '*.go')

# Commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GORUN := $(GOCMD) run
GOLINT := golangci-lint
PROTOC := protoc

# Default Target
all: build

# Build the Go application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/cmd/api.go

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run Producer Script
producer:
	@if [ -z "$(topic)" ]; then \
		echo "Usage: make producer topic=<topic_name>"; \
	else \
		echo "Running producer for topic: $(topic)..."; \
		$(GORUN) $(SCRIPTS_DIR)/producer.go $(topic); \
	fi

# Run Consumer Script
consumer:
	@if [ -z "$(topic)" ]; then \
		echo "Usage: make consumer topic=<topic_name>"; \
	else \
		echo "Running consumer for topic: $(topic)..."; \
		$(GORUN) $(SCRIPTS_DIR)/consumer.go $(topic); \
	fi

# Create a topic
create-topic:
	@echo "Creating topic: $(TOPIC) with strategy: $(STRATEGY)..."
	@go run $(SCRIPTS_DIR)/create_topic.go $(TOPIC) $(STRATEGY)

# List all topics
list-topics:
	@echo "Listing all topics..."
	@go run $(SCRIPTS_DIR)/list_topics.go

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) ./... -v

# Linting
lint:
	@echo "Linting code..."
	$(GOLINT) run ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Generate code (if applicable, e.g., for mockgen, stringer, etc.)
generate:
	@echo "Generating code..."
	$(GOCMD) generate ./...

# Compile Protobuf files
protoc:
	@echo "Compiling Protobuf files..."
	$(PROTOC) -I. -Ipkg/third_party --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) --grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt=logtostderr=true --openapiv2_out=pkg/swagger --openapiv2_opt=logtostderr=true --proto_path=$(PROTO_DIR) $(PROTO_DIR)/*.proto

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor

# Run everything needed for CI/CD
ci: fmt lint test build

.PHONY: all build run test lint fmt generate protoc clean deps ci
