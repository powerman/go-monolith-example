package netx

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/powerman/must"
)

//nolint:gochecknoglobals // By design.
var (
	usedTCPPort   = make(map[int]int)
	usedTCPPortMu sync.Mutex
)

// UnusedTCPPort returns random unique unused TCP port at host.
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

	usedTCPPortMu.Lock()
	used := usedTCPPort[port]
	usedTCPPort[port]++
	usedTCPPortMu.Unlock()
	if used > 0 {
		const maxRecursion = 3
		if used > maxRecursion {
			panic(fmt.Sprintf("same TCP port returned multiple times: %d", port))
		}
		return UnusedTCPPort(host)
	}

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
