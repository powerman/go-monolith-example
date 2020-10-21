// Package config provides configuration shared by microservices.
//
// Default values can be obtained from various sources (constants,
// environment variables, etc.) and then overridden by flags.
//
// As configuration is global you can get it only once for safety:
// you can call only one of Get… functions and call it just once.
package config

import (
	"strconv"

	"github.com/powerman/appcfg"

	"github.com/powerman/go-monolith-example/pkg/def"
)

// EnvPrefix defines common prefix for environment variables.
const EnvPrefix = "MONO_"

// Shared contains configurable values shared by microservices.
type Shared struct {
	ExampleAddrPort        appcfg.Port           `env:"EXAMPLE_ADDR_PORT"`
	ExampleMetricsAddrPort appcfg.Port           `env:"EXAMPLE_METRICS_ADDR_PORT"`
	MonoAddrHost           appcfg.NotEmptyString `env:"MONO_ADDR_HOST"`
	MonoMetricsAddrHost    appcfg.NotEmptyString `env:"MONO_METRICS_ADDR_HOST"`
	XMySQLAddrHost         appcfg.NotEmptyString `env:"X_MYSQL_ADDR_HOST"`
	XMySQLAddrPort         appcfg.Port           `env:"X_MYSQL_ADDR_PORT"`
	XNATSAddrUrls          appcfg.NotEmptyString `env:"X_NATS_ADDR_URLS"`
	XSTANClusterID         appcfg.NotEmptyString `env:"X_STAN_CLUSTER_ID"`
}

// Default ports for metrics.
const (
	MetricsPort = 16000 + iota // XXX Used in ms/metrics/service.go.
	ExampleMetricsPort
)

// Default ports.
const (
	_ = 17000 + iota
	ExamplePort
)

var shared = &Shared{ //nolint:gochecknoglobals // Config is global anyway.
	ExampleAddrPort:        appcfg.MustPort(strconv.Itoa(ExamplePort)),
	ExampleMetricsAddrPort: appcfg.MustPort(strconv.Itoa(ExampleMetricsPort)),
	MonoAddrHost:           appcfg.MustNotEmptyString(def.Hostname),
	MonoMetricsAddrHost:    appcfg.MustNotEmptyString(def.Hostname),
	XMySQLAddrPort:         appcfg.MustPort("3306"),
}

// Get updates config defaults (from env) and returns shared config.
func Get() (*Shared, error) {
	defer cleanup()

	fromEnv := appcfg.NewFromEnv(EnvPrefix)
	err := appcfg.ProvideStruct(shared, fromEnv)
	if err != nil {
		return nil, err
	}
	return shared, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	shared = nil
}
