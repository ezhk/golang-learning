package api

// Generate gRPC gateway.
//go:generate protoc -I . -I ${GOPATH}/src -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go-grpc_out . --go-grpc_opt require_unimplemented_servers=false --go_out . --go_opt paths=source_relative --openapiv2_out . --openapiv2_opt logtostderr=true --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true banner-api.proto
