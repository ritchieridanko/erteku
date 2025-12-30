# === Variables ===
COMPOSE_FILE := compose.yaml
SHARED_DIR := shared
PROTOS_DIR := $(SHARED_DIR)/protos
APIS_DIR := $(SHARED_DIR)/contract/apis
EVENTS_DIR := $(SHARED_DIR)/contract/events

help:
	@echo "Available commands:"
	@echo " make build-proto                  Build the proto files"
	@echo " make delete-proto                 Delete the .pb.go files"
	@echo " make docker-build                 Build the services"
	@echo " make docker-up                    Run the services"
	@echo " make docker-down                  Drop the services"
	@echo " make docker-migrator-up           Apply all up database migrations"
	@echo " make docker-migrator-down         Apply all down database migrations"

# ---------- Proto Commands ----------
build-proto:
	make delete-proto
	protoc \
		--proto_path=$(PROTOS_DIR) \
		--go_out=$(APIS_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(APIS_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTOS_DIR)/*/*_api.proto
	protoc \
		--proto_path=$(PROTOS_DIR) \
		--go_out=$(EVENTS_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(EVENTS_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTOS_DIR)/*/*_event.proto

delete-proto:
	@find $(APIS_DIR) -name "*.pb.go" -type f -delete
	@find $(EVENTS_DIR) -name "*.pb.go" -type f -delete

# ---------- Docker Commands ----------
docker-build:
	docker compose -f $(COMPOSE_FILE) build

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down -v

docker-migrator-up:
	docker compose run --rm auth-migrator -up

docker-migrator-down:
	docker compose run --rm auth-migrator -down 0
