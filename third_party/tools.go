// +build tools

package third_party

import (
	_ "github.com/cheekybits/genny"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/googleapis/api-linter/cmd/api-linter"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/mattn/goveralls"
	_ "github.com/nilslice/protolock/cmd/protolock"
	_ "github.com/powerman/dockerize"
	_ "github.com/yoheimuta/protolint/cmd/protoc-gen-protolint"
	_ "github.com/yoheimuta/protolint/cmd/protolint"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "gotest.tools/gotestsum"
)
