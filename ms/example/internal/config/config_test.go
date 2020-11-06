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
		MySQL: def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr("localhost", 3306),
			User: "example",
			Pass: "",
			DB:   "example",
		}),
		MySQLGooseDir: "ms/example/internal/migrations",
		NATSURLs:      "nats://localhost:4222",
		STANClusterID: "cluster",
		AuthAddrInt:   netx.NewAddr(def.Hostname, config.AuthPortInt),
		Addr:          netx.NewAddr(def.Hostname, config.ExamplePort),
		MetricsAddr:   netx.NewAddr(def.Hostname, config.ExampleMetricsPort),
		Path:          "/rpc",
		TLSCACert:     "ca.crt",
	}

	t.Run("required", func(tt *testing.T) {
		t := check.T(tt)
		require(t, "MySQLPass")
		os.Setenv("MONO__EXAMPLE_MYSQL_AUTH_PASS", "")
	})
	t.Run("default", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe()
		t.Nil(err)
		t.DeepEqual(c, want)
	})
	t.Run("constraint", func(tt *testing.T) {
		t := check.T(tt)
		constraint(t, "MONO__EXAMPLE_MYSQL_AUTH_LOGIN", "", `^MySQLUser .* empty`)
		constraint(t, "MONO__EXAMPLE_MYSQL_DB", "", `^MySQLName .* empty`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("MONO__EXAMPLE_MYSQL_AUTH_LOGIN", "user3")
		os.Setenv("MONO__EXAMPLE_MYSQL_AUTH_PASS", "pass3")
		os.Setenv("MONO__EXAMPLE_MYSQL_DB", "db3")
		c, err := testGetServe()
		t.Nil(err)
		want.MySQL = def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr("localhost", 3306),
			User: "user3",
			Pass: "pass3",
			DB:   "db3",
		})
		t.DeepEqual(c, want)
	})
	t.Run("flag", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe(
			"--mysql.host=mysql4",
			"--mysql.port=43306",
			"--example.mysql.user=user4",
			"--example.mysql.pass=pass4",
			"--example.mysql.dbname=db4",
			"--nats.urls=nats://nats4:4222",
			"--stan.cluster_id=cluster4",
			"--auth.host.int=authhost4",
			"--auth.port.int=44",
			"--host=host4",
			"--host-int=metrics4",
			"--example.port=8004",
			"--example.metrics.port=4",
		)
		t.Nil(err)
		want.MySQL = def.NewMySQLConfig(def.MySQLConfig{
			Addr: netx.NewAddr("mysql4", 43306),
			User: "user4",
			Pass: "pass4",
			DB:   "db4",
		})
		want.NATSURLs = "nats://nats4:4222"
		want.STANClusterID = "cluster4"
		want.AuthAddrInt = netx.NewAddr("authhost4", 44)
		want.Addr = netx.NewAddr("host4", 8004)
		want.MetricsAddr = netx.NewAddr("metrics4", 4)
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
