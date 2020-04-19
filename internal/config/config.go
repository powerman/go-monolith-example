// Package config provides generic configuration shared by microservices.
//
// Default values can be load from different sources (constants,
// environment variables, etc.) and then override with flags.
//
// As configuration is global you can get it just once for safety.
package config

import (
	"github.com/powerman/appcfg"
	"github.com/powerman/go-monolith-example/internal/def"
)

// EnvPrefix defines common prefix for environment variables.
const EnvPrefix = "MONO_"

// Cfg contains config values shared by microservices.
type Cfg struct {
	MySQLHost     appcfg.NotEmptyString `env:"MYSQL_HOST"`
	MySQLPort     appcfg.Port           `env:"MYSQL_PORT"`
	NATSUrls      appcfg.NotEmptyString `env:"NATS_URLS"`
	STANClusterID appcfg.NotEmptyString `env:"STAN_CLUSTER_ID"`
	MetricsHost   appcfg.NotEmptyString `env:"METRICS_HOST"`
	RPCHost       appcfg.NotEmptyString `env:"RPC_HOST"`
}

//nolint:gochecknoglobals // By design.
var cfg = &Cfg{
	MySQLHost:   appcfg.MustNotEmptyString("localhost"),
	MySQLPort:   appcfg.MustPort("3306"),
	MetricsHost: appcfg.MustNotEmptyString(def.Hostname),
	RPCHost:     appcfg.MustNotEmptyString(def.Hostname),
}

// Get loads default configuration and returns generic config.
// You can get config just once, next call will panic.
func Get() (*Cfg, error) {
	defer cleanup()

	fromEnv := appcfg.NewFromEnv(EnvPrefix)
	err := appcfg.ProvideStruct(cfg, fromEnv)
	return cfg, err
}

func cleanup() {
	cfg = nil
}
