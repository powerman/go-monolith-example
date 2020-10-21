package config

import (
	"os"
	"path/filepath"

	"github.com/powerman/must"
	"github.com/spf13/pflag"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

// MustGetServeTest returns config suitable for use in tests.
func MustGetServeTest() *ServeConfig {
	sharedCfg, err := config.Get()
	must.NoErr(err)
	err = Init(sharedCfg, FlagSets{
		Serve:      pflag.NewFlagSet("", pflag.ContinueOnError),
		GooseMySQL: pflag.NewFlagSet("", pflag.ContinueOnError),
	})
	must.NoErr(err)
	cfg, err := GetServe()
	must.NoErr(err)

	cfg.MySQL.Timeout = def.TestTimeout

	const host = "localhost"
	cfg.Addr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.MetricsAddr = netx.NewAddr(host, 0)

	rootDir, err := os.Getwd()
	must.NoErr(err)
	for _, err := os.Stat(filepath.Join(rootDir, "go.mod")); os.IsNotExist(err) && filepath.Dir(rootDir) != rootDir; _, err = os.Stat(filepath.Join(rootDir, "go.mod")) {
		rootDir = filepath.Dir(rootDir)
	}

	cfg.MySQLGooseDir = filepath.Join(rootDir, cfg.MySQLGooseDir)

	return cfg
}
