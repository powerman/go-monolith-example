package jsonrpc2x_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/rpc-codec/jsonrpc2"

	"github.com/powerman/go-monolith-example/pkg/jsonrpc2x"
)

func TestError(tt *testing.T) {
	t := check.T(tt)

	other := jsonrpc2.NewError(41, "other")
	target := jsonrpc2.NewError(42, "default")
	custom := jsonrpc2.NewError(42, "custom")
	customx := jsonrpc2x.NewError(42, "custom")
	t.False(errors.Is(custom, target))
	t.True(errors.Is(customx, target))
	t.True(errors.Is(customx, jsonrpc2x.NewError(42, "extra")))
	t.False(errors.Is(customx, other))
	t.True(errors.Is(fmt.Errorf("wrapped: %w", customx), target))
	t.True(errors.Is(customx, fmt.Errorf("wrapped2: %w", target)))
	t.True(errors.Is(fmt.Errorf("wrapped: %w", customx), fmt.Errorf("wrapped2: %w", target)))
	t.DeepEqual(jsonrpc2.ServerError(customx), custom)
}
