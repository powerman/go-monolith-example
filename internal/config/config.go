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
	AddrHost               appcfg.NotEmptyString `env:"ADDR_HOST"`
	AddrHostInt            appcfg.NotEmptyString `env:"ADDR_HOST_INT"`
	AuthAddrHost           appcfg.NotEmptyString `env:"AUTH_ADDR_HOST"`
	AuthAddrHostInt        appcfg.NotEmptyString `env:"AUTH_ADDR_HOST_INT"`
	AuthAddrPort           appcfg.Port           `env:"AUTH_ADDR_PORT"`
	AuthAddrPortInt        appcfg.Port           `env:"AUTH_ADDR_PORT_INT"`
	AuthGRPCGWAddrPort     appcfg.Port           `env:"AUTH_GRPCGW_ADDR_PORT"`
	AuthMetricsAddrPort    appcfg.Port           `env:"AUTH_METRICS_ADDR_PORT"`
	ExampleAddrPort        appcfg.Port           `env:"EXAMPLE_ADDR_PORT"`
	ExampleMetricsAddrPort appcfg.Port           `env:"EXAMPLE_METRICS_ADDR_PORT"`
	TLSCACert              appcfg.NotEmptyString `env:"TLS_CA_CERT"`
	XMySQLAddrHost         appcfg.NotEmptyString `env:"X_MYSQL_ADDR_HOST"`
	XMySQLAddrPort         appcfg.Port           `env:"X_MYSQL_ADDR_PORT"`
	XNATSAddrUrls          appcfg.NotEmptyString `env:"X_NATS_ADDR_URLS"`
	XPostgresAddrHost      appcfg.NotEmptyString `env:"X_POSTGRES_ADDR_HOST"`
	XPostgresAddrPort      appcfg.Port           `env:"X_POSTGRES_ADDR_PORT"`
	XPostgresDBName        appcfg.NotEmptyString `env:"X_POSTGRES_DB_NAME"`
	XSTANClusterID         appcfg.NotEmptyString `env:"X_STAN_CLUSTER_ID"`
}

// Default ports.
const (
	MonoPort = 17000 + iota
	ExamplePort
	ExampleMetricsPort
	AuthPort
	AuthPortInt
	AuthMetricsPort
	AuthGRPCGWPort
)

var shared = &Shared{ //nolint:gochecknoglobals // Config is global anyway.
	AddrHost:               appcfg.MustNotEmptyString(def.Hostname),
	AddrHostInt:            appcfg.MustNotEmptyString(def.Hostname),
	AuthAddrHost:           appcfg.MustNotEmptyString(def.Hostname),
	AuthAddrHostInt:        appcfg.MustNotEmptyString(def.Hostname),
	AuthAddrPort:           appcfg.MustPort(strconv.Itoa(AuthPort)),
	AuthAddrPortInt:        appcfg.MustPort(strconv.Itoa(AuthPortInt)),
	AuthGRPCGWAddrPort:     appcfg.MustPort(strconv.Itoa(AuthGRPCGWPort)),
	AuthMetricsAddrPort:    appcfg.MustPort(strconv.Itoa(AuthMetricsPort)),
	ExampleAddrPort:        appcfg.MustPort(strconv.Itoa(ExamplePort)),
	ExampleMetricsAddrPort: appcfg.MustPort(strconv.Itoa(ExampleMetricsPort)),
	XMySQLAddrPort:         appcfg.MustPort("3306"),
	XPostgresAddrPort:      appcfg.MustPort("5432"),
	XPostgresDBName:        appcfg.MustNotEmptyString("postgres"),
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
