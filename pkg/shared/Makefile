.PHONY: proto clean help

proto:
	@echo "Generating proto files..."
	protoc --proto_path=proto \
		--go_out=. --go_opt=module=github.com/Ernestgio/Hangout-Planner/pkg/shared \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Ernestgio/Hangout-Planner/pkg/shared \
		proto/file/*.proto
	@echo "Proto files generated successfully!"

clean:
	@echo "Cleaning generated proto files..."
	rm -rf proto/gen/go/file/*.pb.go
	@echo "Clean complete!"

help:
	@echo "Available targets:"
	@echo "  proto  - Generate Go code from proto files"
	@echo "  clean  - Remove generated proto files"
	@echo "  help   - Show this help message"
