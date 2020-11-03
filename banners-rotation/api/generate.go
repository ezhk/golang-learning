package api

// Generate gRPC gateway.
//go:generate protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go-grpc_out ../internal/server --go-grpc_opt require_unimplemented_servers=false --go_out ../internal/server --go_opt paths=source_relative --openapiv2_out ../internal/server --openapiv2_opt logtostderr=true --grpc-gateway_out ../internal/server --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true banner-api.proto
