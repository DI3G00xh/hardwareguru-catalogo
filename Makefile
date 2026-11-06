.PHONY: proto test run seed

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
	go test -v -cover ./internal/...