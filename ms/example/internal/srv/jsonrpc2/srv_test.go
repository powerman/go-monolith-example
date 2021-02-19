package jsonrpc2_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/pkg/def"
)

func fetch(t *check.C, url string, headers ...string) *http.Response {
	t.Helper()
	c := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
	t.Cleanup(cancel)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	t.Must(t.Nil(err))
	for i := 0; i < len(headers); i += 2 {
		req.Header.Add(headers[i], headers[i+1])
	}
	resp, err := c.Do(req)
	t.Must(t.Nil(err))
	return resp
}

func TestCORS(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, url, _ := testNew(t)

	resp := fetch(t, url, "Origin", "google.com")
	t.Equal(resp.StatusCode, 405)
	t.Equal(resp.Header.Get("Access-Control-Allow-Origin"), "*")
}
