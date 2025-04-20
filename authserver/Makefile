.PHONY: build run clean docker docker-run

APP_NAME = authserver
DOCKER_IMAGE = wewe-rss/authserver

all: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) ./cmd/server

run:
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/server

clean:
	@echo "Cleaning..."
	@rm -f $(APP_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

docker:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make clean      - Remove binary files"
	@echo "  make test       - Run tests"
	@echo "  make docker     - Build Docker image"
	@echo "  make docker-run - Run Docker container" 