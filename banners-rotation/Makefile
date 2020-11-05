SERVER_SRC_PATH = ./internal/server
STORAGE_SRC_PATH = ./internal/storage

DOCUMENT_PATH = ./docs

build: clean
	go build -o build/banner .

clean:
	rm build/banner || true

generate:
	# generate gRPC server and REST API proxy
	protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis	\
		--openapiv2_out ${DOCUMENT_PATH} --openapiv2_opt logtostderr=true										\
		--go-grpc_out ${SERVER_SRC_PATH} --go-grpc_opt require_unimplemented_servers=false						\
		--go_out ${SERVER_SRC_PATH}																				\
		--grpc-gateway_out ${SERVER_SRC_PATH} --grpc-gateway_opt logtostderr=true								\
			--grpc-gateway_opt generate_unbound_methods=true													\
		api/banner-api.proto

	# generage server methods
	go run ${SERVER_SRC_PATH}/generate-server-methods/... -object banner -file ${SERVER_SRC_PATH}/generated_banner.go
	go run ${SERVER_SRC_PATH}/generate-server-methods/... -object slot -file ${SERVER_SRC_PATH}/generated_slot.go
	go run ${SERVER_SRC_PATH}/generate-server-methods/... -object group -file ${SERVER_SRC_PATH}/generated_group.go

	# generate storage methods and tests
	go run ${STORAGE_SRC_PATH}/generate-table-methods/... -table banner -file ${STORAGE_SRC_PATH}/generated_banner.go
	go run ${STORAGE_SRC_PATH}/generate-table-methods/... -table slot -file ${STORAGE_SRC_PATH}/generated_slot.go
	go run ${STORAGE_SRC_PATH}/generate-table-methods/... -table group -file ${STORAGE_SRC_PATH}/generated_group.go
	go run ${STORAGE_SRC_PATH}/generate-table-tests/... -table banner -file ${STORAGE_SRC_PATH}/generated_banner_test.go
	go run ${STORAGE_SRC_PATH}/generate-table-tests/... -table slot -file ${STORAGE_SRC_PATH}/generated_slot_test.go
	go run ${STORAGE_SRC_PATH}/generate-table-tests/... -table group -file ${STORAGE_SRC_PATH}/generated_group_test.go

test:
	go test -race ./...

integration-test:
	# run tests separate to avoid superposition
	go test -race -count 5 -tags integration ./internal/storage/...
	go test -race -count 5 -tags integration ./internal/server/...
	go test -race -count 5 -tags integration ./internal/queue/...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --disable exhaustivestruct --disable wrapcheck ./...

run:
	go run . --config=configs/banner.yaml $(filter-out $@,$(MAKECMDGOALS))

grpc-client:
	evans  --path . --path ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --port 9090 internal/api/banner-api.proto

.PHONY: build
