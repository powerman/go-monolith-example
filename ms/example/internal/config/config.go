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

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/appcfg"
	"github.com/spf13/pflag"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
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
	MySQLUser appcfg.NotEmptyString `env:"MYSQL_AUTH_LOGIN"`
	MySQLPass appcfg.String         `env:"MYSQL_AUTH_PASS"`
	MySQLName appcfg.NotEmptyString `env:"MYSQL_DB"`
	GooseDir  appcfg.NotEmptyString
}{ // Defaults, if any:
	MySQLUser: appcfg.MustNotEmptyString(app.ServiceName),
	MySQLName: appcfg.MustNotEmptyString(app.ServiceName),
	GooseDir:  appcfg.MustNotEmptyString(fmt.Sprintf("ms/%s/internal/migrations", app.ServiceName)),
}

// FlagSets for all CLI subcommands which use flags to set config values.
type FlagSets struct {
	Serve      *pflag.FlagSet
	GooseMySQL *pflag.FlagSet
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

	appcfg.AddPFlag(fs.GooseMySQL, &shared.XMySQLAddrHost, "mysql.host", "host to connect to MySQL")
	appcfg.AddPFlag(fs.GooseMySQL, &shared.XMySQLAddrPort, "mysql.port", "port to connect to MySQL")
	appcfg.AddPFlag(fs.GooseMySQL, &own.MySQLUser, "mysql.user", "MySQL username")
	appcfg.AddPFlag(fs.GooseMySQL, &own.MySQLPass, "mysql.pass", "MySQL password")
	appcfg.AddPFlag(fs.GooseMySQL, &own.MySQLName, "mysql.dbname", "MySQL database name")

	pfx := app.ServiceName + "."
	appcfg.AddPFlag(fs.Serve, &shared.XMySQLAddrHost, "mysql.host", "host to connect to MySQL")
	appcfg.AddPFlag(fs.Serve, &shared.XMySQLAddrPort, "mysql.port", "port to connect to MySQL")
	appcfg.AddPFlag(fs.Serve, &own.MySQLUser, pfx+"mysql.user", "MySQL username")
	appcfg.AddPFlag(fs.Serve, &own.MySQLPass, pfx+"mysql.pass", "MySQL password")
	appcfg.AddPFlag(fs.Serve, &own.MySQLName, pfx+"mysql.dbname", "MySQL database name")
	appcfg.AddPFlag(fs.Serve, &shared.XNATSAddrUrls, "nats.urls", "URLs to connect to NATS (separated by comma)")
	appcfg.AddPFlag(fs.Serve, &shared.XSTANClusterID, "stan.cluster_id", "STAN cluster ID")
	appcfg.AddPFlag(fs.Serve, &shared.AddrHost, "host", "host to serve")
	appcfg.AddPFlag(fs.Serve, &shared.AddrHostInt, "host-int", "internal host to serve")
	appcfg.AddPFlag(fs.Serve, &shared.ExampleAddrPort, pfx+"port", "port to serve")
	appcfg.AddPFlag(fs.Serve, &shared.ExampleMetricsAddrPort, pfx+"metrics.port", "port to serve Prometheus metrics")

	return nil
}

// ServeConfig contains configuration for subcommand.
type ServeConfig struct {
	MySQL         *mysql.Config
	MySQLGooseDir string
	NATSURLs      string
	STANClusterID string
	Addr          netx.Addr
	MetricsAddr   netx.Addr
}

// GetServe validates and returns configuration for subcommand.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		MySQL: def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr(shared.XMySQLAddrHost.Value(&err), shared.XMySQLAddrPort.Value(&err)),
			User: own.MySQLUser.Value(&err),
			Pass: own.MySQLPass.Value(&err),
			DB:   own.MySQLName.Value(&err),
		}),
		MySQLGooseDir: own.GooseDir.Value(&err),
		NATSURLs:      shared.XNATSAddrUrls.Value(&err),
		STANClusterID: shared.XSTANClusterID.Value(&err),
		Addr:          netx.NewAddr(shared.AddrHost.Value(&err), shared.ExampleAddrPort.Value(&err)),
		MetricsAddr:   netx.NewAddr(shared.AddrHostInt.Value(&err), shared.ExampleMetricsAddrPort.Value(&err)),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, own, shared)
	}
	return c, nil
}

func GetGooseMySQL() (c *cobrax.GooseMySQLConfig, err error) {
	defer cleanup()

	c = &cobrax.GooseMySQLConfig{
		MySQL: def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr(shared.XMySQLAddrHost.Value(&err), shared.XMySQLAddrPort.Value(&err)),
			User: own.MySQLUser.Value(&err),
			Pass: own.MySQLPass.Value(&err),
			DB:   own.MySQLName.Value(&err),
		}),
		MySQLGooseDir: own.GooseDir.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.GooseMySQL, own, shared)
	}
	return c, nil
}

// Cleanup must be called by all Get* functions to ensure second call to
// any of them will panic.
func cleanup() {
	own = nil
	shared = nil
}
