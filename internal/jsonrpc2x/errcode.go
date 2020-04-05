package jsonrpc2x

import (
	"errors"
	"strconv"

	jsonrpc2pkg "github.com/powerman/rpc-codec/jsonrpc2"
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
	return strconv.Itoa(err.(*jsonrpc2pkg.Error).Code)
}

func dropcode(err error) error {
	if err2 := new(jsonrpc2pkg.Error); errors.As(err, &err2) {
		return errors.New(err2.Message)
	}
	return err
}
