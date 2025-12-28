# === Variables ===
BINARY_DIR := bin
MIGRATOR_BIN := $(BINARY_DIR)/migrator

help:
	@echo "Available commands:"
	@echo " make setup-database               Create the database"
	@echo " make drop-database                Drop the database"
	@echo " make run-migrator-up              Apply all up migrations"
	@echo " make run-migrator-down            Apply N down migrations"
	@echo " make run-migrator-down-all        Apply all down migrations"
	@echo " make build-migrator               Build the migrator"
	@echo " make build-and-run-migrator-up    Build and run the migrator up"

# ---------- Database Commands ----------
setup-database:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'erteku_auth_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE erteku_auth_db"

drop-database:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS erteku_auth_db"

# ---------- Migrator Commands ----------
run-migrator-up:
	go run cmd/migrator/main.go -up

run-migrator-down:
	go run cmd/migrator/main.go -down $(steps)

run-migrator-down-all:
	go run cmd/migrator/main.go -down 0

build-migrator:
	go build -o $(MIGRATOR_BIN) cmd/migrator/main.go

build-and-run-migrator-up:
	make build-migrator
	./$(MIGRATOR_BIN) -up
