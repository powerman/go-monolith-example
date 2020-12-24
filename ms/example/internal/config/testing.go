package config

import (
	"os"
	"path/filepath"
	"strings"

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

	const hostInt = "127.0.0.1"
	const host = "localhost"
	cfg.AuthAddrInt = netx.NewAddr(hostInt, 0) // Invalid value for easier bug detection if not changed.
	cfg.BindAddr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.BindMetricsAddr = netx.NewAddr(hostInt, 0)

	rootDir, err := os.Getwd()
	must.NoErr(err)
	for _, err := os.Stat(filepath.Join(rootDir, "go.mod")); os.IsNotExist(err) && filepath.Dir(rootDir) != rootDir; _, err = os.Stat(filepath.Join(rootDir, "go.mod")) {
		rootDir = filepath.Dir(rootDir)
	}

	for _, path := range []*string{
		&cfg.GooseMySQLDir,
		&cfg.TLSCACert,
	} {
		if !strings.HasPrefix(*path, "/") {
			*path = filepath.Join(rootDir, *path)
		}
	}

	return cfg
}
