.PHONY: build run clean test tidy migrate-up migrate-down migrate-create migrate-force migrate-version db-up db-down db-logs

# Variables
MIGRATIONS_DIR ?= internal/db/migrations

# Build the application
build:
	go build -o bin/app ./cmd/app

# Run the application
run: build
	./bin/app

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Tidy dependencies
tidy:
	go mod tidy

format:
	go fmt ./...

# ==================== Migrations ====================
# Requires: migrate CLI (https://github.com/golang-migrate/migrate)
# Install: brew install golang-migrate (macOS)
#          or: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run all up migrations
migrate-up:
	migrate -database "$${DB_URL}" -path $(MIGRATIONS_DIR) up

# Run N up migrations (usage: make migrate-up-n N=1)
migrate-up-n:
	migrate -database "$${DB_URL}" -path $(MIGRATIONS_DIR) up $(N)

# Rollback all migrations
migrate-down:
	migrate -database "$${DB_URL}" -path $(MIGRATIONS_DIR) down

# Rollback N migrations (usage: make migrate-down-n N=1)
migrate-down-n:
	migrate -database "$${DB_URL}" -path $(MIGRATIONS_DIR) down $(N)

# Force set migration version (usage: make migrate-force V=1)
migrate-force:
	migrate -database "$${DB_URL}" -path $(MIGRATIONS_DIR) force $(V)

# Show current migration version
migrate-version:
	migrate -database "$${DB_URL}" -path $(MIGRATIONS_DIR) version

# Create a new migration (usage: make migrate-create NAME=create_users_table)
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
