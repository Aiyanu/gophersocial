include .env
# Database migration configuration
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: test
test:
	@go test -v ./...

# Create a new migration file with the given name
# Usage: make migration name=create_users_table
.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

# Apply all pending migrations
.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

# Revert migrations
# Usage: make migrate-down [N] where N is number of migrations to revert
.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

# Force set the database version when in a dirty state
# Usage: make migrate-force version=X
.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))

# Show current migration version and status
.PHONY: migrate-version
migrate-version:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) version

# Seed the database with initial data
.PHONY: seed
seed: 
	@go run cmd/migrate/seed/main.go

# Generate API documentation
.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
