package config

import (
	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/must"
	"github.com/spf13/pflag"
)

// MustGetTest provides a convenient way to get config for tests.
// It'll call Init for you. You can get config just once,
// next call to this or other method which returns config will panic.
func MustGetTest() (cfg *ServeConfig) {
	genericCfg, err := config.Get()
	if err == nil {
		err = Init(genericCfg, FlagSets{
			Serve: pflag.NewFlagSet("", pflag.ContinueOnError),
			Goose: pflag.NewFlagSet("", pflag.ContinueOnError),
		})
	}
	if err == nil {
		cfg, err = GetServe()
	}
	must.NoErr(err)

	var connectTimeout = 3 * def.TestSecond
	cfg.MySQLConfig.Timeout = connectTimeout

	return cfg
}
