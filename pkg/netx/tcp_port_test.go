package netx_test

import (
	"sort"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/pkg/netx"
)

func TestUnusedTCPPort(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	ports := make([]int, 10)
	for i := 0; i < len(ports); i++ {
		ports[i] = netx.UnusedTCPPort("127.0.0.1")
	}
	sort.Ints(ports)
	for i := 1; i < len(ports); i++ {
		t.NotEqual(ports[i-1], ports[i])
	}
}
