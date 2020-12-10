package config

import (
	"os"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/pqx"

	"github.com/powerman/go-monolith-example/internal/config"
	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

func Test(t *testing.T) {
	want := &ServeConfig{
		Postgres: def.NewPostgresConfig(pqx.Config{
			Host:        "postgres",
			Port:        5432,
			DBName:      "postgres",
			User:        "auth",
			Pass:        "authpass",
			SSLRootCert: "ca.crt",
		}),
		GoosePostgresDir: "ms/auth/internal/migrations",
		Addr:             netx.NewAddr(def.Hostname, config.AuthPort),
		AddrInt:          netx.NewAddr(def.Hostname, config.AuthPortInt),
		GRPCGWAddr:       netx.NewAddr(def.Hostname, config.AuthGRPCGWPort),
		MetricsAddr:      netx.NewAddr(def.Hostname, config.AuthMetricsPort),
		Secret:           []byte("s3cr3t"),
		TLSCACert:        "ca.crt",
		TLSCert:          "tls.crt",
		TLSCertInt:       "tls-int.crt",
		TLSKey:           "tls.key",
		TLSKeyInt:        "tls-int.key",
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
		require(t, "PostgresPass")
		os.Setenv("MONO__AUTH_POSTGRES_AUTH_PASS", "authpass")
	})
	t.Run("default", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe()
		t.Nil(err)
		t.DeepEqual(c, want)
	})
	t.Run("constraint", func(tt *testing.T) {
		t := check.T(tt)
		constraint(t, "MONO__AUTH_POSTGRES_AUTH_LOGIN", "", `^PostgresUser .* empty`)
		constraint(t, "MONO__AUTH_POSTGRES_AUTH_PASS", "", `^PostgresPass .* empty`)
		constraint(t, "MONO__AUTH_SECRET", "", `^Secret .* empty`)
		constraint(t, "MONO__AUTH_TLS_CERT", "", `^TLSCert .* empty`)
		constraint(t, "MONO__AUTH_TLS_CERT_INT", "", `^TLSCertInt .* empty`)
		constraint(t, "MONO__AUTH_TLS_KEY", "", `^TLSKey .* empty`)
		constraint(t, "MONO__AUTH_TLS_KEY_INT", "", `^TLSKeyInt .* empty`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("MONO__AUTH_POSTGRES_AUTH_LOGIN", "auth3")
		os.Setenv("MONO__AUTH_POSTGRES_AUTH_PASS", "authpass3")
		os.Setenv("MONO__AUTH_SECRET", "secret3")
		os.Setenv("MONO__AUTH_TLS_CERT", "tls3.crt")
		os.Setenv("MONO__AUTH_TLS_CERT_INT", "tls3-int.crt")
		os.Setenv("MONO__AUTH_TLS_KEY", "tls3.key")
		os.Setenv("MONO__AUTH_TLS_KEY_INT", "tls3-int.key")
		c, err := testGetServe()
		t.Nil(err)
		want.Postgres.User = "auth3"
		want.Postgres.Pass = "authpass3"
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
			"--postgres.host=localhost4",
			"--postgres.port=45432",
			"--postgres.dbname=postgres4",
			"--auth.postgres.user=auth4",
			"--auth.postgres.pass=authpass4",
			"--host=host4",
			"--host-int=hostint4",
			"--auth.port=8004",
			"--auth.port.int=9004",
			"--auth.grpcgw.port=7004",
			"--auth.metrics.port=4",
			"--auth.secret=secret4\u212B\u0041\u030A\u00C5", // From https://www.unicode.org/reports/tr15/#Singletons_Figure.
		)
		t.Nil(err)
		want.Postgres.Host = "localhost4"
		want.Postgres.Port = 45432
		want.Postgres.DBName = "postgres4"
		want.Postgres.User = "auth4"
		want.Postgres.Pass = "authpass4"
		want.Addr = netx.NewAddr("host4", 8004)
		want.AddrInt = netx.NewAddr("hostint4", 9004)
		want.GRPCGWAddr = netx.NewAddr("host4", 7004)
		want.MetricsAddr = netx.NewAddr("hostint4", 4)
		want.Secret = []byte("secret4\u0041\u030A\u0041\u030A\u0041\u030A") // NFD form.
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
