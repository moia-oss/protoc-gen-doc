//go:build tools
// +build tools

package tools

import (
	_ "github.com/goreleaser/goreleaser/v2"
	_ "github.com/haya14busa/goverage"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
