PROTO_DIR := proto
OUT_DIR := internal/embedding
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

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
