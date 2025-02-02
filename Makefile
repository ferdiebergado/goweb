# Project Name (adjust as needed)
PROJECT_NAME = goweb

# Binary Name
BINARY_NAME = $(PROJECT_NAME)

# Go Modules Path (if using Go Modules)
GO_MODULE_PATH = ./...  # Or specify specific packages like: ./api/... ./web/...

# Build Directory
BUILD_DIR = build

# Versioning (you can automate this later)
VERSION = v0.1.0

# Go Flags (e.g., for race detector)
GO_FLAGS = -race

.PHONY: $(wildcard *)

%:
	@true

# Build Targets
## default: Show usage information
default:
	@echo "Usage:"
	@sed -n 's/^## //p' Makefile | column -t -s ':' --table-columns TARGET," DESCRIPTION"," EXAMPLE"

## build: Build the project
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GO_FLAGS) -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(GO_MODULE_PATH)
	@echo "Build complete!"

run: build
	@echo "Running $(BINARY_NAME) $(VERSION)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test $(GO_FLAGS) $(GO_MODULE_PATH) -coverprofile=coverage.out

test-cover: test
	@go tool cover -html=coverage.out

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@echo "Clean complete!"

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(PROJECT_NAME):$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(PROJECT_NAME):$(VERSION)  # Adjust port mapping

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run -v $(GO_MODULE_PATH) # Make sure golangci-lint.yml is configured

format:
	@echo "Running go fmt..."
	@go fmt $(GO_MODULE_PATH)

migrate-up: # Example for database migrations (using migrate tool)
	@echo "Running database migrations (up)..."
	@migrate -path ./migrations -database "postgres://..." -steps 1 up # Replace with your DB connection string

migrate-down:
	@echo "Running database migrations (down)..."
	@migrate -path ./migrations -database "postgres://..." down # Replace with your DB connection string

## gen: Generate source files
gen:
	@command -v mockgen >/dev/null || go install go.uber.org/mock/mockgen@latest
	go generate -v ./...

## tidy: Add missing or remove unused modules
tidy:
	go mod tidy

dev:
	@GO_FLAGS=
	@ENV=development
	@$(MAKE) run

prod:
	@GO_FLAGS=-ldflags="-s -w" # Example production flags
	@ENV=production
	@$(MAKE) run
