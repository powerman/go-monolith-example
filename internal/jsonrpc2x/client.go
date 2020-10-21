package jsonrpc2x

import (
	"errors"
	"io"
	"net/rpc"
	"sync"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

// Client provides an easier way to use jsonrpc2.Client.
type Client struct {
	sync.Mutex
	*jsonrpc2.Client
	url string
}

// NewHTTPClient wraps jsonrpc2.NewHTTPClient.
func NewHTTPClient(url string) *Client {
	return &Client{
		Client: jsonrpc2.NewHTTPClient(url),
		url:    url,
	}
}

// Call invokes the named function, waits for it to complete, and returns its error status.
// It also applies jsonrpc2.WrapError to returned error and automatically
// handles rpc.ErrShutdown and io.ErrUnexpectedEOF.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	err := jsonrpc2.WrapError(c.Client.Call(serviceMethod, args, reply))
	if errors.Is(err, rpc.ErrShutdown) || errors.Is(err, io.ErrUnexpectedEOF) {
		c.Lock()
		defer c.Unlock()
		c.Client = jsonrpc2.NewHTTPClient(c.url)
	}
	return err
}
