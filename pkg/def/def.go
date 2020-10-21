// Package def provides default values for both commands and tests.
package def

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/powerman/getenv"
	"github.com/powerman/must"
	"github.com/powerman/sqlxx"
	"github.com/prometheus/client_golang/prometheus"
)

func init() { //nolint:gochecknoinits // Ensure time.Now() assigned to global vars uses UTC.
	// Make time.Now()==time.Now().UTC() https://github.com/golang/go/issues/19486
	time.Local = nil
}

// Init must be called once before using this package.
// It provides common initialization for both commands and tests.
func Init() error {
	// Make sure no one occasionally uses global objects.
	http.DefaultServeMux = nil
	prometheus.DefaultRegisterer = nil
	prometheus.DefaultGatherer = nil

	must.AbortIf = must.PanicIf

	sqlx.NameMapper = sqlxx.ToSnake

	setupLog()

	if hostnameErr != nil {
		return fmt.Errorf("os.Hostname: %w", hostnameErr)
	}
	return getenv.LastErr()
}
