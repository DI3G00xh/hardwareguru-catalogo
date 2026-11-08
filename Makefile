.PHONY: proto test test-coverage test-integration run seed

proto:
	protoc --proto_path=proto/hardware \
	       --go_out=proto/hardware --go_opt=paths=source_relative \
	       --go-grpc_out=proto/hardware --go-grpc_opt=paths=source_relative \
	       proto/hardware/hardware.proto

run:
	go run cmd/server/main.go

seed:
	go run cmd/seeder/main.go

test:
	go test -v -cover ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

test-integration:
	go test -v ./tests/...
