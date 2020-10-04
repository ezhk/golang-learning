build:
	go build -o ./build ./...

test:
	go test -race ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

run:
	go run . --config=configs/calendar.yaml $(filter-out $@,$(MAKECMDGOALS))

# brew tap ktr0731/evans && brew install evans
grpc-client:
	evans  --path . --path ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --port 9090 internal/server/grpc/internalgrpc.proto

# DSN="user=postgres password=postgres dbname=postgres sslmode=disable" make migrate
migrate:
	goose -dir migrations postgres "${DSN}" up

generate:
	go generate ./...

.PHONY: build
