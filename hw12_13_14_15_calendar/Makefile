build: clean test
	go build -o ./build/calendar -tags=server,sender,scheduler .

clean:
	find build -not -name "README.md" -delete || true

test:
	go test -race ./...

integration-test:
	go test -tags integration ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --build-tags server,sender,scheduler --disable exhaustivestruct ./...

run:
	go run -tags=server,sender,scheduler . --config=configs/calendar.yaml $(filter-out $@,$(MAKECMDGOALS))

# brew tap ktr0731/evans && brew install evans
grpc-client:
	evans  --path . --path ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --port 9090 internal/server/grpc/internalgrpc.proto

# DSN="user=postgres password=postgres dbname=postgres sslmode=disable" make migrate
migrate:
	goose -dir migrations postgres "${DSN}" up

generate:
	go generate ./...

up:
	docker-compose -f docker/docker-compose.yaml up --scale bdd-tests=0

down:
	docker-compose -f docker/docker-compose.yaml down

bdd:
	docker-compose -f docker/docker-compose.yaml run --rm bdd-tests

.PHONY: build
