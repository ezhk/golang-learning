build:
	go build -o ./build ./...

test:
	go test -race ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

run:
	go run . --config=conf/calendar.yaml $(filter-out $@,$(MAKECMDGOALS))

# DSN="user=postgres password=postgres dbname=postgres sslmode=disable" make migrate
migrate:
	goose -dir migrations postgres "${DSN}" up

generate:
	go generate ./...

.PHONY: build
