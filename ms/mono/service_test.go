// +build integration

package mono

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)

	s := &Service{}

	const host = "localhost"
	s.cfg.BindAddr = netx.NewAddr(host, netx.UnusedTCPPort(host))

	ctxStartup, cancel := context.WithTimeout(ctx, def.TestTimeout)
	defer cancel()
	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error)
	go func() { errc <- s.RunServe(ctxStartup, ctxShutdown, shutdown) }()
	defer func() {
		shutdown()
		t.Nil(<-errc, "RunServe")
	}()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, s.cfg.BindAddr), "connect to HTTP service"))

	client := gorequest.New().
		Timeout(def.TestTimeout).
		Retry(30, def.TestSecond/10, http.StatusServiceUnavailable)
	endpoint := fmt.Sprintf("http://%s/", s.cfg.BindAddr)

	{ // health-check
		resp, body, errs := client.Clone().Get(endpoint + "health-check").End()
		t.Nil(errs)
		t.Equal(resp.StatusCode, 200)
		t.Equal(body, "OK")
	}
}
