// Package api describes microservice auth's gRPC API.
package api

//go:generate bash -e -o pipefail -c "rm -f *.pb.*; d=./$(git rev-parse --show-prefix); cd $(git rev-parse --show-toplevel); PATH=$(go list -tags=tools -f '{{range .Imports}}{{println .}}{{end}}' . | grep protoc-gen | xargs gobin -m -p | xargs dirname | sed -z -e 's,\\n,:,g')$DOLLAR{PATH} gobin -m -run github.com/bufbuild/buf/cmd/buf generate --template $DOLLAR{d}buf.gen.yaml $(find $DOLLAR{d} -maxdepth 1 -name '*.proto' -printf '--file %p ')"
