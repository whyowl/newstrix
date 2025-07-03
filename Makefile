PROTO_DIR := proto
OUT_DIR := internal/embedding
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

DB_URL := postgres://news:password@localhost:5432/newsdb?sslmode=disable

generate:
	@echo "Generating Go gRPC code from .proto files..."
	@protoc \
		--go_out=$(OUT_DIR) \
		--go-grpc_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Done."

clean:
	@echo "Cleaning generated files..."
	@rm -f $(OUT_DIR)/proto/*.pb.go
	@echo "Cleaned."

migrate-up:
	@docker run --rm \
		-v ${CURDIR}/internal/db/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		up

migrate-down:
	@docker run --rm \
		-v ${CURDIR}/internal/db/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		down

create-migration:
	@migrate create -ext sql -dir internal/db/migrations -seq $(name)

api:
	@go run ./cmd/api

fetcher:
	@go run ./cmd/fetcher

# --- DEV MODE ---
dev-up:
	@echo "Run dev-environment..."
	@docker compose -f docker-compose.dev.yml up -d --build

dev-down:
	@echo "Stop dev-environment..."
	@docker compose -f docker-compose.dev.yml down

dev:
	@$(MAKE) dev-up
	@$(MAKE) migrate-up
	@echo "Run API and fetcher local..."
	@go run ./cmd/fetcher #& \
	#go run ./cmd/api

# --- PROD MODE ---
prod-up:
	@echo "Run whole project in Docker (prod)..."
	@docker compose -f docker-compose.yml up -d --build

prod-down:
	@echo "Stop prod-environment..."
	@docker compose -f docker-compose.yml down