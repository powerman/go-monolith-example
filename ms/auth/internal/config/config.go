// Package config provides configurations for subcommands.
//
// It consists of both configuration values shared by all
// microservices and values specific to this microservice.
//
// Default values can be obtained from various sources (constants,
// environment variables, etc.) and then overridden by flags.
//
// As configuration is global you can get it only once for safety:
// you can call only one of Get… functions and call it just once.
package config

import (
	"fmt"
	"strings"

	"github.com/powerman/appcfg"
	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/norm"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

// EnvPrefix defines common prefix for environment variables.
var envPrefix = fmt.Sprintf("%s_%s_", config.EnvPrefix, strings.ToUpper(app.ServiceName)) //nolint:gochecknoglobals // Const.

type SharedCfg = config.Shared

var shared *SharedCfg //nolint:gochecknoglobals // Config is global anyway.

// Own configurable values of the microservice.
//
// If microservice may runs in different ways (e.g. using CLI subcommands)
// then these subcommands may use subset of these values.
var own = &struct { //nolint:gochecknoglobals // Config is global anyway.
	Secret     appcfg.NotEmptyString `env:"SECRET"`
	TLSCert    appcfg.NotEmptyString `env:"TLS_CERT"`
	TLSCertInt appcfg.NotEmptyString `env:"TLS_CERT_INT"`
	TLSKey     appcfg.NotEmptyString `env:"TLS_KEY"`
	TLSKeyInt  appcfg.NotEmptyString `env:"TLS_KEY_INT"`
}{ // Defaults, if any:
}

// FlagSets for all CLI subcommands which use flags to set config values.
type FlagSets struct {
	Serve *pflag.FlagSet
}

var fs FlagSets //nolint:gochecknoglobals // Flags are global anyway.

// Init updates config defaults (from env) and setup subcommands flags.
//
// Init must be called once before using this package.
func Init(sharedCfg *SharedCfg, flagsets FlagSets) error {
	shared, fs = sharedCfg, flagsets

	fromEnv := appcfg.NewFromEnv(envPrefix)
	err := appcfg.ProvideStruct(own, fromEnv)
	if err != nil {
		return err
	}

	pfx := app.ServiceName + "."
	appcfg.AddPFlag(fs.Serve, &shared.AddrHost, "host", "host to serve")
	appcfg.AddPFlag(fs.Serve, &shared.AddrHostInt, "host-int", "internal host to serve")
	appcfg.AddPFlag(fs.Serve, &shared.AuthAddrPort, pfx+"port", "port to serve")
	appcfg.AddPFlag(fs.Serve, &shared.AuthAddrPortInt, pfx+"port.int", "port to serve internal API")
	appcfg.AddPFlag(fs.Serve, &shared.AuthGRPCGWAddrPort, pfx+"grpcgw.port", "port to serve grpc-gateway")
	appcfg.AddPFlag(fs.Serve, &shared.AuthMetricsAddrPort, pfx+"metrics.port", "port to serve Prometheus metrics")
	appcfg.AddPFlag(fs.Serve, &own.Secret, pfx+"secret", "secret used for hashing passwords")

	return nil
}

// ServeConfig contains configuration for subcommand.
type ServeConfig struct {
	Addr        netx.Addr
	AddrInt     netx.Addr
	GRPCGWAddr  netx.Addr
	MetricsAddr netx.Addr
	Secret      []byte
	TLSCACert   string
	TLSCert     string
	TLSCertInt  string
	TLSKey      string
	TLSKeyInt   string
}

// GetServe validates and returns configuration for subcommand.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		Addr:        netx.NewAddr(shared.AddrHost.Value(&err), shared.AuthAddrPort.Value(&err)),
		AddrInt:     netx.NewAddr(shared.AddrHostInt.Value(&err), shared.AuthAddrPortInt.Value(&err)),
		GRPCGWAddr:  netx.NewAddr(shared.AddrHost.Value(&err), shared.AuthGRPCGWAddrPort.Value(&err)),
		MetricsAddr: netx.NewAddr(shared.AddrHostInt.Value(&err), shared.AuthMetricsAddrPort.Value(&err)),
		Secret:      norm.NFD.Bytes([]byte(own.Secret.Value(&err))),
		TLSCACert:   shared.TLSCACert.Value(&err),
		TLSCert:     own.TLSCert.Value(&err),
		TLSCertInt:  own.TLSCertInt.Value(&err),
		TLSKey:      own.TLSKey.Value(&err),
		TLSKeyInt:   own.TLSKeyInt.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, own, shared)
	}
	return c, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	own = nil
	shared = nil
}
