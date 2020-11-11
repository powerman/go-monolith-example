package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/powerman/must"
	"github.com/spf13/pflag"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

// MustGetServeTest returns config suitable for use in tests.
func MustGetServeTest() *ServeConfig {
	sharedCfg, err := config.Get()
	must.NoErr(err)
	err = Init(sharedCfg, FlagSets{
		Serve: pflag.NewFlagSet("", pflag.ContinueOnError),
	})
	must.NoErr(err)
	cfg, err := GetServe()
	must.NoErr(err)

	const hostInt = "127.0.0.1"
	const host = "localhost"
	cfg.Addr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.AddrInt = netx.NewAddr(hostInt, netx.UnusedTCPPort(hostInt))
	cfg.GRPCGWAddr = netx.NewAddr(host, netx.UnusedTCPPort(host))
	cfg.MetricsAddr = netx.NewAddr(hostInt, 0)

	rootDir, err := os.Getwd()
	must.NoErr(err)
	for _, err := os.Stat(filepath.Join(rootDir, "go.mod")); os.IsNotExist(err) && filepath.Dir(rootDir) != rootDir; _, err = os.Stat(filepath.Join(rootDir, "go.mod")) {
		rootDir = filepath.Dir(rootDir)
	}

	for _, path := range []*string{
		&cfg.TLSCACert,
		&cfg.TLSCert,
		&cfg.TLSCertInt,
		&cfg.TLSKey,
		&cfg.TLSKeyInt,
	} {
		if !strings.HasPrefix(*path, "/") {
			*path = filepath.Join(rootDir, *path)
		}
	}

	return cfg
}
