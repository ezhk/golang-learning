build:
	go build -o ./build ./...

test:
	go test -race ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

run:
	go run . --config=conf/calendar.yaml $(filter-out $@,$(MAKECMDGOALS))

# DSN="user=postgres password=postgres dbname=postgres sslmode=disable" make migrate
migrate:
	goose -dir migrations postgres "${DSN}" up

generate:
	protoc --go_out=internal/server/grpc internal/server/grpc/internalgrpc.proto

.PHONY: build
