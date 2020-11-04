// +build tools

package tools

import (
	_ "github.com/cheekybits/genny"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/mattn/goveralls"
	_ "github.com/powerman/dockerize"
	_ "gotest.tools/gotestsum"
)
