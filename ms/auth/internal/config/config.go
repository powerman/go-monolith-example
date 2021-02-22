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
	"github.com/powerman/pqx"
	"github.com/powerman/sensitive"
	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/norm"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/ms/auth/internal/app"
	"github.com/powerman/go-monolith-example/pkg/cobrax"
	"github.com/powerman/go-monolith-example/pkg/def"
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
	PostgresUser     appcfg.NotEmptyString `env:"POSTGRES_AUTH_LOGIN"`
	PostgresPass     appcfg.NotEmptyString `env:"POSTGRES_AUTH_PASS"`
	GoosePostgresDir appcfg.NotEmptyString
	Secret           appcfg.NotEmptyString `env:"SECRET"`
	TLSCert          appcfg.NotEmptyString `env:"TLS_CERT"`
	TLSCertInt       appcfg.NotEmptyString `env:"TLS_CERT_INT"`
	TLSKey           appcfg.NotEmptyString `env:"TLS_KEY"`
	TLSKeyInt        appcfg.NotEmptyString `env:"TLS_KEY_INT"`
}{ // Defaults, if any:
	PostgresUser:     appcfg.MustNotEmptyString(app.ServiceName),
	GoosePostgresDir: appcfg.MustNotEmptyString(fmt.Sprintf("ms/%s/internal/migrations", app.ServiceName)),
}

// FlagSets for all CLI subcommands which use flags to set config values.
type FlagSets struct {
	Serve         *pflag.FlagSet
	GoosePostgres *pflag.FlagSet
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

	appcfg.AddPFlag(fs.GoosePostgres, &shared.XPostgresAddrHost, "postgres.host", "host to connect to PostgreSQL")
	appcfg.AddPFlag(fs.GoosePostgres, &shared.XPostgresAddrPort, "postgres.port", "port to connect to PostgreSQL")
	appcfg.AddPFlag(fs.GoosePostgres, &shared.XPostgresDBName, "postgres.dbname", "PostgreSQL database name")
	appcfg.AddPFlag(fs.GoosePostgres, &own.PostgresUser, "postgres.user", "PostgreSQL username")
	appcfg.AddPFlag(fs.GoosePostgres, &own.PostgresPass, "postgres.pass", "PostgreSQL password")

	pfx := app.ServiceName + "."
	appcfg.AddPFlag(fs.Serve, &shared.XPostgresAddrHost, "postgres.host", "host to connect to PostgreSQL")
	appcfg.AddPFlag(fs.Serve, &shared.XPostgresAddrPort, "postgres.port", "port to connect to PostgreSQL")
	appcfg.AddPFlag(fs.Serve, &shared.XPostgresDBName, "postgres.dbname", "PostgreSQL database name")
	appcfg.AddPFlag(fs.Serve, &own.PostgresUser, pfx+"postgres.user", "PostgreSQL username")
	appcfg.AddPFlag(fs.Serve, &own.PostgresPass, pfx+"postgres.pass", "PostgreSQL password")
	appcfg.AddPFlag(fs.Serve, &shared.AddrHost, "host", "host to serve")
	appcfg.AddPFlag(fs.Serve, &shared.AddrHostInt, "host-int", "internal host to serve")
	appcfg.AddPFlag(fs.Serve, &shared.AuthAddrHost, "auth.host", "ms/auth API host")
	appcfg.AddPFlag(fs.Serve, &shared.AuthAddrPort, "auth.port", "ms/auth API port")
	appcfg.AddPFlag(fs.Serve, &shared.AuthAddrPortInt, "auth.port-int", "ms/auth internal API port")
	appcfg.AddPFlag(fs.Serve, &shared.AuthGRPCGWAddrPort, "auth.grpcgw.port", "ms/auth OpenAPI port")
	appcfg.AddPFlag(fs.Serve, &shared.AuthMetricsAddrPort, "auth.metrics.port", "ms/auth Prometheus metrics port")
	appcfg.AddPFlag(fs.Serve, &own.Secret, pfx+"secret", "secret used for hashing passwords")

	return nil
}

// ServeConfig contains configuration for subcommand.
type ServeConfig struct {
	Postgres         *def.PostgresConfig
	GoosePostgresDir string
	AuthAddr         netx.Addr
	BindAddr         netx.Addr
	BindAddrInt      netx.Addr
	BindGRPCGWAddr   netx.Addr
	BindMetricsAddr  netx.Addr
	Secret           sensitive.Bytes
	TLSCACert        string
	TLSCert          string
	TLSCertInt       string
	TLSKey           string
	TLSKeyInt        string
}

// GetServe validates and returns configuration for subcommand.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		Postgres: def.NewPostgresConfig(pqx.Config{
			Host:        shared.XPostgresAddrHost.Value(&err),
			Port:        shared.XPostgresAddrPort.Value(&err),
			DBName:      shared.XPostgresDBName.Value(&err),
			User:        own.PostgresUser.Value(&err),
			Pass:        own.PostgresPass.Value(&err),
			SSLRootCert: shared.TLSCACert.Value(&err),
		}),
		GoosePostgresDir: own.GoosePostgresDir.Value(&err),
		AuthAddr:         netx.NewAddr(shared.AuthAddrHost.Value(&err), shared.AuthAddrPort.Value(&err)),
		BindAddr:         netx.NewAddr(shared.AddrHost.Value(&err), shared.AuthAddrPort.Value(&err)),
		BindAddrInt:      netx.NewAddr(shared.AddrHostInt.Value(&err), shared.AuthAddrPortInt.Value(&err)),
		BindGRPCGWAddr:   netx.NewAddr(shared.AddrHost.Value(&err), shared.AuthGRPCGWAddrPort.Value(&err)),
		BindMetricsAddr:  netx.NewAddr(shared.AddrHostInt.Value(&err), shared.AuthMetricsAddrPort.Value(&err)),
		Secret:           norm.NFD.Bytes([]byte(own.Secret.Value(&err))),
		TLSCACert:        shared.TLSCACert.Value(&err),
		TLSCert:          own.TLSCert.Value(&err),
		TLSCertInt:       own.TLSCertInt.Value(&err),
		TLSKey:           own.TLSKey.Value(&err),
		TLSKeyInt:        own.TLSKeyInt.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, own, shared)
	}
	return c, nil
}

func GetGoosePostgres() (c *cobrax.GoosePostgresConfig, err error) {
	defer cleanup()

	c = &cobrax.GoosePostgresConfig{
		Postgres: def.NewPostgresConfig(pqx.Config{
			Host:        shared.XPostgresAddrHost.Value(&err),
			Port:        shared.XPostgresAddrPort.Value(&err),
			DBName:      shared.XPostgresDBName.Value(&err),
			User:        own.PostgresUser.Value(&err),
			Pass:        own.PostgresPass.Value(&err),
			SSLRootCert: shared.TLSCACert.Value(&err),
		}),
		GoosePostgresDir: own.GoosePostgresDir.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.GoosePostgres, own, shared)
	}
	return c, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	own = nil
	shared = nil
}
