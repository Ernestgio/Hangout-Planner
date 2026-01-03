.PHONY: proto clean help

proto:
	@echo "Generating proto files..."
	protoc --proto_path=proto \
		--go_out=proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
		proto/file/file_service.proto proto/file/file_messages.proto
	@echo "Proto files generated successfully!"

clean:
	@echo "Cleaning generated proto files..."
	rm -rf proto/gen/file/*.pb.go
	@echo "Clean complete!"
