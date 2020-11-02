build: clean
	go build -o build/banner .

clean:
	rm build/banner || true

generate:
	go generate ./...

test:
	# run tests separate to avoid superposition
	go test -race -count 5 -tags integration ./internal/storage/...
	go test -race -count 5 -tags integration ./internal/server/...
	go test -race -count 5 -tags integration ./internal/queue/...
	go test -count 1 ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --disable exhaustivestruct --disable wrapcheck ./...

run:
	go run . --config=config/banner.yaml $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
