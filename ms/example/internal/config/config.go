// Package config provides microservice configuration.
//
// It consists of both generic configuration values shared by all
// microservices and values specific to this microservice.
//
// Default values can be load from different sources (constants,
// environment variables, etc.) and then override with flags.
//
// As configuration is global you can get it just once for safety.
package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/appcfg"
	"github.com/powerman/go-monolith-example/internal/cobrax"
	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/netx"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/spf13/pflag"
)

//nolint:gochecknoglobals // Const.
var envPrefix = fmt.Sprintf("%s%s_", config.EnvPrefix, strings.ToUpper(app.ServiceName))

// GenericCfg is a synonym for convenience.
type GenericCfg = config.Cfg

var generic *GenericCfg //nolint:gochecknoglobals // By design.

type svcCfg struct {
	MySQLUser   appcfg.NotEmptyString `env:"DB_USER"`
	MySQLPass   appcfg.String         `env:"DB_PASS"`
	MySQLName   appcfg.NotEmptyString `env:"DB_NAME"`
	GooseDir    appcfg.NotEmptyString
	MetricsPort appcfg.Port `env:"METRICS_PORT"`
	RPCPort     appcfg.Port `env:"RPC_PORT"`
}

var svc = &svcCfg{ //nolint:gochecknoglobals // By design.
	MySQLUser:   appcfg.MustNotEmptyString(app.ServiceName),
	MySQLName:   appcfg.MustNotEmptyString(app.ServiceName),
	GooseDir:    appcfg.MustNotEmptyString(fmt.Sprintf("ms/%s/internal/migrations", app.ServiceName)),
	MetricsPort: appcfg.MustPort(strconv.Itoa(def.ExampleMetricsPort)),
	RPCPort:     appcfg.MustPort(strconv.Itoa(def.ExampleRPCPort)),
}

// FlagSets contains all flagsets which needs flags tied to service config.
type FlagSets struct{ Serve, Goose *pflag.FlagSet }

var fs FlagSets //nolint:gochecknoglobals // By design.

// Init loads default configuration and setup flags to override defaults.
//
// Init must be called once before using this package.
func Init(genericCfg *GenericCfg, flagsets FlagSets) error {
	generic, fs = genericCfg, flagsets

	fromEnv := appcfg.NewFromEnv(envPrefix)
	err := appcfg.ProvideStruct(svc, fromEnv)
	if err != nil {
		return err
	}

	config.AddFlag(fs.Goose, &generic.MySQLHost, "mysql.host", "host to connect to MySQL")
	config.AddFlag(fs.Goose, &generic.MySQLPort, "mysql.port", "port to connect to MySQL")
	config.AddFlag(fs.Goose, &svc.MySQLUser, "mysql.user", "MySQL username")
	config.AddFlag(fs.Goose, &svc.MySQLPass, "mysql.pass", "MySQL password")
	config.AddFlag(fs.Goose, &svc.MySQLName, "mysql.dbname", "MySQL database name")

	pfx := app.ServiceName + "."
	config.AddFlag(fs.Serve, &generic.MySQLHost, "mysql.host", "host to connect to MySQL")
	config.AddFlag(fs.Serve, &generic.MySQLPort, "mysql.port", "port to connect to MySQL")
	config.AddFlag(fs.Serve, &svc.MySQLUser, pfx+"mysql.user", "MySQL username")
	config.AddFlag(fs.Serve, &svc.MySQLPass, pfx+"mysql.pass", "MySQL password")
	config.AddFlag(fs.Serve, &svc.MySQLName, pfx+"mysql.dbname", "MySQL database name")
	config.AddFlag(fs.Serve, &generic.NATSUrls, "nats.urls", "URLs to connect to NATS (separated by comma)")
	config.AddFlag(fs.Serve, &generic.STANClusterID, "stan.cluster_id", "STAN cluster ID")
	config.AddFlag(fs.Serve, &generic.MetricsHost, "metrics.host", "host to serve Prometheus metrics")
	config.AddFlag(fs.Serve, &svc.MetricsPort, pfx+"metrics.port", "port to serve Prometheus metrics")
	config.AddFlag(fs.Serve, &generic.RPCHost, "rpc.host", "host to serve JSON-RPC 2.0")
	config.AddFlag(fs.Serve, &svc.RPCPort, pfx+"rpc.port", "port to serve JSON-RPC 2.0")

	return nil
}

// GetGoose validates and returns config. You can get config just once,
// next call to this or other method which returns config will panic.
func GetGoose() (c *cobrax.GooseConfig, err error) {
	defer cleanup()

	c = &cobrax.GooseConfig{
		MySQLConfig: def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr(generic.MySQLHost.Value(&err), generic.MySQLPort.Value(&err)),
			User: svc.MySQLUser.Value(&err),
			Pass: svc.MySQLPass.Value(&err),
			DB:   svc.MySQLName.Value(&err),
		}),
		GooseDir: svc.GooseDir.Value(&err),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Goose, svc, generic)
	}
	return c, nil
}

// ServeConfig contain configuration for serve command.
type ServeConfig struct {
	NATSUrls      string
	STANClusterID string
	MySQLConfig   *mysql.Config
	GooseDir      string
	MetricsAddr   netx.Addr
	RPCAddr       netx.Addr
}

// GetServe validates and returns config. You can get config just once,
// next call to this or other method which returns config will panic.
func GetServe() (c *ServeConfig, err error) {
	defer cleanup()

	c = &ServeConfig{
		NATSUrls:      generic.NATSUrls.Value(&err),
		STANClusterID: generic.STANClusterID.Value(&err),
		MySQLConfig: def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr(generic.MySQLHost.Value(&err), generic.MySQLPort.Value(&err)),
			User: svc.MySQLUser.Value(&err),
			Pass: svc.MySQLPass.Value(&err),
			DB:   svc.MySQLName.Value(&err),
		}),
		GooseDir:    svc.GooseDir.Value(&err),
		MetricsAddr: netx.NewAddr(generic.MetricsHost.Value(&err), svc.MetricsPort.Value(&err)),
		RPCAddr:     netx.NewAddr(generic.RPCHost.Value(&err), svc.RPCPort.Value(&err)),
	}
	if err != nil {
		return nil, appcfg.WrapPErr(err, fs.Serve, svc, generic)
	}
	return c, nil
}

func cleanup() {
	svc = nil
	generic = nil
}
