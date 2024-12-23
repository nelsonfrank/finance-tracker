include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed: 
	@go run cmd/migrate/seed/main.go

.PHONY: build
build:
	@echo "Building..."

	@go build -o ./bin/main ./cmd/api

.PHONY: watch
watch:
	@echo "Checking if air is installed..."
	@if command -v air > /dev/null; then \
		echo "Air found. Running air..."; \
		air; \
	else \
		echo "Air is not installed."; \
		read -p "Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			echo "Installing air..."; \
			go install github.com/air-verse/air@latest; \
			echo "Running air..."; \
			air; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi
