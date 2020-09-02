build:
	go build -o ./build ./...

test:
	go test -race ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

run:
	go run . --config=conf/calendar.yaml $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
