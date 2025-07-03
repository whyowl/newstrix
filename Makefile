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

up:
	@echo "Starting containers..."
	@docker compose up -d --build

down:
	@echo "Stopping containers..."
	@docker compose down

migrate-up:
	@docker run --rm \
		-v $(shell pwd)/internal/db/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		up

migrate-down:
	@docker run --rm \
		-v $(shell pwd)/internal/db/migrations:/migrations \
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

dev: up migrate-up
	@echo "Starting API and Fetcher..."
	@go run ./cmd/api & \
	go run ./cmd/fetcher