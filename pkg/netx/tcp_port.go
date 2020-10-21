package netx

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/powerman/must"
)

// UnusedTCPPort returns random unused TCP port at host.
func UnusedTCPPort(host string) (port int) {
	var portStr string
	ln, err := net.Listen("tcp", host+":0")
	if err == nil {
		err = ln.Close()
	}
	if err == nil {
		_, portStr, err = net.SplitHostPort(ln.Addr().String())
	}
	if err == nil {
		port, err = strconv.Atoi(portStr)
	}
	must.NoErr(err)
	return port
}

// WaitTCPPort tries to connect to addr until success or ctx.Done.
func WaitTCPPort(ctx context.Context, addr fmt.Stringer) error {
	const delay = time.Second / 20
	var dialer net.Dialer
	for ; ctx.Err() == nil; time.Sleep(delay) {
		conn, err := dialer.DialContext(ctx, "tcp", addr.String())
		if err == nil {
			return conn.Close()
		}
	}
	return ctx.Err()
}
