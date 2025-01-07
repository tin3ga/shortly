# Variables
APP_NAME := shortly
BUILD_DIR := ./bin

GOOSE_DRIVER := postgres
GOOSE_DBSTRING := db_url # replace with connection string
GOOSE_MIGRATION_DIR := sql/schema/


.PHONY: install
install: ## Create a install dev dependencies
	@echo "Installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Installing sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Installing air..."
	@go install github.com/air-verse/air@latest


.PHONY: new-migration
new-migration: ## Create a new migration
	@read -p "Enter migration name: " name; \
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" create "$$name" sql

.PHONY: migrate-up
migrate-up: ## Apply all up migrations
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down



.PHONY: run-dev
run-dev:## run air - Live reload
	@air


.PHONY: sqlc
sqlc:## Generate fully type-safe idiomatic code from SQL(folder:sql/queries)
	@sqlc generate



.PHONY: gen-docs
gen-docs:## Generate and format swagger docs
	@swag init && swag fmt



.PHONY: build
build:## Build executable
	@echo "Building the project..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

.PHONY: run
run:## Run executable
	@echo "Running the project..."
	@$(BUILD_DIR)/$(APP_NAME)

.PHONY: clean
clean:## Clean project
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR) tmp *.exe


.PHONY: help
help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

