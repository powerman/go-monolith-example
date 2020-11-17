// Package api describes microservice auth's gRPC API.
package api

// - `buf` must be run in project's root
// - protoc plugins must be in $PATH, so use `gobin` to set $PATH
// - remove useless swagger.json needlessly generated for internal API
// - set swagger.json mtime to latest of it's sources, to ensure statik won't update it needlessly
//go:generate bash -e -o pipefail -c "rm -f *.pb.* *.swagger.json; d=./$(git rev-parse --show-prefix); cd $(git rev-parse --show-toplevel); PATH=$(go list -tags=tools -f '{{range .Imports}}{{println .}}{{end}}' . | grep protoc-gen | xargs gobin -m -p | xargs dirname | sed -z -e 's,\\n,:,g')$DOLLAR{PATH} gobin -m -run github.com/bufbuild/buf/cmd/buf generate --template $DOLLAR{d}buf.gen.yaml --path $DOLLAR{d}; rm -f $DOLLAR{d}*_int.swagger.json; touch -r $(ls -t $DOLLAR{d}*.proto $DOLLAR{d}*.openapi.yml | head -n 1) $DOLLAR{d}service.swagger.json"
