package jsonrpc2x

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

func codes(errs []error) []string {
	codes := make([]string, len(errs))
	for i := range errs {
		codes[i] = code(errs[i])
	}
	return codes
}

func code(err error) string {
	if err == nil {
		return ""
	}
	if rpcerr := new(jsonrpc2.Error); errors.As(err, &rpcerr) {
		return strconv.Itoa(rpcerr.Code)
	}
	panic(fmt.Sprintf("not a jsonrpc2.Error: %#+v", err))
}

func dropcode(err error) error {
	if err2 := new(jsonrpc2.Error); errors.As(err, &err2) {
		return errors.New(err2.Message) //nolint:goerr113 // By design.
	}
	return err
}
