.PHONY: build run test clean

# Build the application
build:
	go build -o pos-api.exe .

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f pos-api.exe

# Install dependencies
deps:
	go mod tidy

# Create database (requires psql)
createdb:
	psql -U postgres -c "CREATE DATABASE pos_db;"

# Drop database (requires psql)
dropdb:
	psql -U postgres -c "DROP DATABASE IF EXISTS pos_db;"

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run
