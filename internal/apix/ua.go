//go:generate -command mockgen sh -c "$(git rev-parse --show-toplevel)/.gobincache/$DOLLAR{DOLLAR}0 \"$DOLLAR{DOLLAR}@\"" mockgen
//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE

package apix

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/powerman/must"
	"github.com/powerman/structlog"

	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

const (
	defaultTimeout = 30 * time.Second
	maxBodySize    = 1 * 1024 * 1024
)

var errRespTooLarge = errors.New("HTTP response is too large")

// UserAgent is a convenience wrapper for http.Client, suitable for
// fetching small responses (because it reads full response in memory).
type UserAgent interface {
	// Do sends an HTTP request and returns an HTTP response.
	// It returns body for convenience - it's same as can be read from resp.Body.
	// It saves body to file if log level is Debug and cfg.Debug and cfg.DumpDir is set.
	// Returned resp.Body doesn't needs to be closed.
	Do(ctx Ctx, req *http.Request, skip int) (_ *http.Response, body []byte, _ error)
	// Log resp (both HTTP request and response).
	// It does nothing if cfg.Debug is false.
	Log(ctx Ctx, resp *http.Response, body []byte)
}

// UserAgentConfig contains configuration for UserAgent.
type UserAgentConfig struct {
	Timeout     time.Duration // Default: 30s.
	MaxBodySize int           // Default: 1MB.
	Debug       bool          // Log response.
	DumpDir     string        // If not empty and Debug - dump response body to files in this dir.
}

type userAgent struct {
	cfg    UserAgentConfig
	client *http.Client // Default: &http.Client{}. Feel free to change as needed.
}

// NewUserAgent creates and returns new UserAgent.
func NewUserAgent(cfg UserAgentConfig) UserAgent {
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.MaxBodySize == 0 {
		cfg.MaxBodySize = maxBodySize
	}
	return &userAgent{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (x *userAgent) Do(ctx Ctx, req *http.Request, skip int) (_ *http.Response, body []byte, _ error) {
	log := structlog.FromContext(ctx, nil)
	ctx, cancel := context.WithTimeout(ctx, x.cfg.Timeout)
	defer cancel()

	resp, err := x.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	defer log.WarnIfFail(resp.Body.Close)

	switch {
	case resp.ContentLength > maxBodySize:
		return nil, nil, errRespTooLarge
	case resp.ContentLength < 0:
		resp.Body = ioutil.NopCloser(io.LimitReader(resp.Body, maxBodySize+1))
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if len(body) > maxBodySize {
		return nil, nil, errRespTooLarge
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	x.dump(ctx, resp, body, skip+1)
	return resp, body, nil
}

func (x *userAgent) Log(ctx Ctx, resp *http.Response, body []byte) {
	log := structlog.FromContext(ctx, nil)
	if !(log.IsDebug() && x.cfg.Debug) {
		return
	}

	const maxLogBodyBytes = 1019 // Prime number to increase chance last line won't be full and cut mark will be easier to spot.
	const cutMark = "....."
	if len(body) > maxLogBodyBytes {
		part := append(body[:maxLogBodyBytes:maxLogBodyBytes], []byte(cutMark)...) //nolint:gocritic // Not same slice.
		defer func(r io.ReadCloser) { resp.Body = r }(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewReader(part))
	}

	var dump bytes.Buffer
	must.NoErr(resp.Request.Write(&dump))
	must.NoErr(dump.WriteByte('\n'))
	must.NoErr(resp.Write(&dump))
	log.Debug("response", "dump", dump.String())
}

func (x *userAgent) dump(ctx Ctx, resp *http.Response, body []byte, skip int) {
	log := structlog.FromContext(ctx, nil)
	if !(log.IsDebug() && x.cfg.Debug && x.cfg.DumpDir != "") {
		return
	}

	dumpName := reflectx.CallerPkg(skip+1) + "." + reflectx.CallerMethodName(skip+1)
	if ext, _ := mime.ExtensionsByType(resp.Header.Get("Content-Type")); len(ext) > 0 {
		dumpName += ext[0]
	} else {
		dumpName += ".data"
	}
	dumpPath := filepath.Join(x.cfg.DumpDir, dumpName)
	err := ioutil.WriteFile(dumpPath, body, 0o600)
	if err != nil {
		log.Warn("failed to save response body", "err", err)
		return
	}
	log.Debug("saved response body", "file", dumpPath)
}
