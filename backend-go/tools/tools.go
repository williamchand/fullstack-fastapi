//go:build tools
// +build tools

package tools // import "ignore.tools"

import (
	// Protocol buffer tools
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-lint"

	// gRPC and protobuf
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	// gRPC Gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"

	// SQLC
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"

	// Database migrations
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"

	// Linting
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"

	// Formatiing
	_ "mvdan.cc/gofumpt"
)
