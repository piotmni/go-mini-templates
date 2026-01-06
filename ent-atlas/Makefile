.PHONY: generate migrate-diff migrate-apply migrate-status

# Database connection
DB_URL ?= postgres://pastebin:pastebin@localhost:5432/pastebin?sslmode=disable

# Generate Ent code
generate:
	go generate ./ent

# Create a new versioned migration
# Usage: make migrate-diff name=<migration_name>
migrate-diff:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-diff name=<migration_name>"; exit 1; fi
	atlas migrate diff $(name) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://postgres/15/test?search_path=public"

# Apply pending migrations to the database
migrate-apply:
	atlas migrate apply \
		--dir "file://ent/migrate/migrations" \
		--url "$(DB_URL)"

# Show migration status
migrate-status:
	atlas migrate status \
		--dir "file://ent/migrate/migrations" \
		--url "$(DB_URL)"

# Lint migration files
migrate-lint:
	atlas migrate lint \
		--dev-url "docker://postgres/15/test?search_path=public" \
		--dir "file://ent/migrate/migrations" \
		--latest 1


# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run:
	DB_URL="$(DB_URL)" go run ./cmd/api
