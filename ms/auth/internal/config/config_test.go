package config

import (
	"os"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

func Test(t *testing.T) {
	want := &ServeConfig{
		Addr:        netx.NewAddr(def.Hostname, config.AuthPort),
		AddrInt:     netx.NewAddr(def.Hostname, config.AuthPortInt),
		GRPCGWAddr:  netx.NewAddr(def.Hostname, config.AuthGRPCGWPort),
		MetricsAddr: netx.NewAddr(def.Hostname, config.AuthMetricsPort),
		Secret:      []byte("s3cr3t"),
		TLSCACert:   "ca.crt",
		TLSCert:     "tls.crt",
		TLSCertInt:  "tls-int.crt",
		TLSKey:      "tls.key",
		TLSKeyInt:   "tls-int.key",
	}

	t.Run("required", func(tt *testing.T) {
		t := check.T(tt)
		require(t, "TLSKeyInt")
		os.Setenv("MONO__AUTH_TLS_KEY_INT", "tls-int.key")
		require(t, "TLSKey")
		os.Setenv("MONO__AUTH_TLS_KEY", "tls.key")
		require(t, "TLSCertInt")
		os.Setenv("MONO__AUTH_TLS_CERT_INT", "tls-int.crt")
		require(t, "TLSCert")
		os.Setenv("MONO__AUTH_TLS_CERT", "tls.crt")
		require(t, "Secret")
		os.Setenv("MONO__AUTH_SECRET", "s3cr3t")
	})
	t.Run("default", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe()
		t.Nil(err)
		t.DeepEqual(c, want)
	})
	t.Run("constraint", func(tt *testing.T) {
		t := check.T(tt)
		constraint(t, "MONO__AUTH_SECRET", "", `^Secret .* empty`)
		constraint(t, "MONO__AUTH_TLS_CERT", "", `^TLSCert .* empty`)
		constraint(t, "MONO__AUTH_TLS_CERT_INT", "", `^TLSCertInt .* empty`)
		constraint(t, "MONO__AUTH_TLS_KEY", "", `^TLSKey .* empty`)
		constraint(t, "MONO__AUTH_TLS_KEY_INT", "", `^TLSKeyInt .* empty`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("MONO__AUTH_SECRET", "secret3")
		os.Setenv("MONO__AUTH_TLS_CERT", "tls3.crt")
		os.Setenv("MONO__AUTH_TLS_CERT_INT", "tls3-int.crt")
		os.Setenv("MONO__AUTH_TLS_KEY", "tls3.key")
		os.Setenv("MONO__AUTH_TLS_KEY_INT", "tls3-int.key")
		c, err := testGetServe()
		t.Nil(err)
		want.Secret = []byte("secret3")
		want.TLSCert = "tls3.crt"
		want.TLSCertInt = "tls3-int.crt"
		want.TLSKey = "tls3.key"
		want.TLSKeyInt = "tls3-int.key"
		t.DeepEqual(c, want)
	})
	t.Run("flag", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe(
			"--host=host4",
			"--host-int=hostint4",
			"--auth.port=8004",
			"--auth.port.int=9004",
			"--auth.grpcgw.port=7004",
			"--auth.metrics.port=4",
			"--auth.secret=secret4", // TODO Test norm.NFD.
		)
		t.Nil(err)
		want.Addr = netx.NewAddr("host4", 8004)
		want.AddrInt = netx.NewAddr("hostint4", 9004)
		want.GRPCGWAddr = netx.NewAddr("host4", 7004)
		want.MetricsAddr = netx.NewAddr("hostint4", 4)
		want.Secret = []byte("secret4")
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
