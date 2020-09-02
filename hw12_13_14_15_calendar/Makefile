build:
	go build -o ./build ./...

test:
	go test -race ./...

lint:
	golangci-lint run ./...

run:
	go run . --config=conf/calendar.yaml $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
