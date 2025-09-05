# --- Configuration ---
# Variables are defined here for easy modification and reuse.
CONTAINER_NAME := go-task-db
DATAMODEL_DIR := ./datamodel
DATASET_DIR := ./dataset
MYSQL_USER := root
MYSQL_HOST := 127.0.0.1
MYSQL_PORT := 4000
SCHEMAS := tasking
SCRIPT_NAME := *.sql

# --- Targets ---

# Default target when you just run `make`
default: help

# Starts the TiDB container if not already running
db-up:
	@if [ -z "$$(docker ps -q -f name=$(CONTAINER_NAME))" ]; then \
		echo "--- Starting $(CONTAINER_NAME) container ---"; \
		docker run --rm -d \
			--name $(CONTAINER_NAME) \
			-v $(PWD)/datamodel:/opt/datamodel \
			-v $(PWD)/dataset:/opt/dataset \
			-p $(MYSQL_PORT):4000 \
			pingcap/tidb:v8.5.2; \
		docker exec -it $(CONTAINER_NAME) yum install mariadb -y; \
	else \
		echo "--- Database container is already running ---"; \
	fi

# Stops the database container.
db-down:
	@echo "--- Stopping $(CONTAINER_NAME) container ---"
	@docker stop $(CONTAINER_NAME) || true

# Runs SQL migrations against the database.
db-migrate:
	@echo "--- Starting migrations ---"
	@for dbs in $(SCHEMAS); do \
		echo ">> Creating schema $$dbs"; \
		docker exec -i $(CONTAINER_NAME) sh -c "mysql --host 127.0.0.1 --port $(MYSQL_PORT) -u root < /opt/datamodel/$$dbs-schema-create.sql"; \
		\
		echo ">> Running datamodel scripts for $$dbs"; \
		for s in $$(find $(DATAMODEL_DIR)/$$dbs -name '$(SCRIPT_NAME)' -exec basename {} \; | sort); do \
			echo "   Running: $$s"; \
			docker exec -i $(CONTAINER_NAME) sh -c "mysql --host 127.0.0.1 -D $$dbs --port $(MYSQL_PORT) -u root < /opt/datamodel/$$dbs/$$s"; \
		done; \
		\
		echo ">> Running dataset scripts for $$dbs"; \
		for s in $$(find $(DATASET_DIR)/$$dbs -name '$(SCRIPT_NAME)' -exec basename {} \; | sort); do \
			echo "   Running: $$s"; \
			docker exec -i $(CONTAINER_NAME) sh -c "mysql --host 127.0.0.1 -D $$dbs --port $(MYSQL_PORT) -u root < /opt/dataset/$$dbs/$$s"; \
		done; \
	done
	@echo "✅ Migrations finished successfully."

# A convenient workflow target to completely reset the database.
db-reset: db-down db-up db-migrate
	@echo "✅ Database reset complete."

# Watch for changes and automatically restart the application.
watch:
	air -c ./.air.conf

# Build the application.
build:
	go build -o ./go-task ./app/cmd/go-task

# Test
test:
	go test -v ./app/controller/handler/... ./app/model/...

# A simple help target to explain how to use the Makefile.
help:
	@echo "Available commands:"
	@echo "  make db-up       - Starts the TiDB container if not running."
	@echo "  make db-down     - Stops the TiDB container."
	@echo "  make db-migrate  - Runs all .sql files from datamodel/ and dataset/."
	@echo "  make db-reset    - Stops, starts, and migrates the database."
	@echo "  make help        - Shows this help message."
	@echo "  make watch       - Starts the application with live reload using air."
	@echo "  make build       - Builds the application binary."

# --- Housekeeping ---
.PHONY: default help db-up db-down db-migrate db-reset watch build test