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
			Addr:   netx.NewAddr("localhost", 3306),
			User:   "example",
			Pass:   "",
			DBName: "example",
		}),
		GooseMySQLDir:   "ms/example/internal/migrations",
		NATSURLs:        "nats://localhost:4222",
		STANClusterID:   "cluster",
		AuthAddrInt:     netx.NewAddr(def.Hostname, config.AuthPortInt),
		BindAddr:        netx.NewAddr(def.Hostname, config.ExamplePort),
		BindMetricsAddr: netx.NewAddr(def.Hostname, config.ExampleMetricsPort),
		Path:            "/rpc",
		TLSCACert:       "ca.crt",
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
		constraint(t, "MONO__EXAMPLE_MYSQL_DB_NAME", "", `^MySQLDBName .* empty`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("MONO__EXAMPLE_MYSQL_AUTH_LOGIN", "user3")
		os.Setenv("MONO__EXAMPLE_MYSQL_AUTH_PASS", "pass3")
		os.Setenv("MONO__EXAMPLE_MYSQL_DB_NAME", "db3")
		c, err := testGetServe()
		t.Nil(err)
		want.MySQL = def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr("localhost", 3306),
			User:   "user3",
			Pass:   "pass3",
			DBName: "db3",
		})
		t.DeepEqual(c, want)
	})
	t.Run("flag", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe(
			"--mysql.host=localhost4",
			"--mysql.port=4200",
			"--example.mysql.user=user4",
			"--example.mysql.pass=pass4",
			"--example.mysql.dbname=db4",
			"--nats.urls=nats://nats4:4222",
			"--stan.cluster_id=cluster4",
			"--host=host4",
			"--host-int=hostint4",
			"--auth.host-int=authhost4int",
			"--auth.port-int=4104",
			"--example.port=4101",
			"--example.metrics.port=4102",
		)
		t.Nil(err)
		want.MySQL = def.NewMySQLConfig(def.MySQLConfig{
			Addr:   netx.NewAddr("localhost4", 4200),
			User:   "user4",
			Pass:   "pass4",
			DBName: "db4",
		})
		want.NATSURLs = "nats://nats4:4222"
		want.STANClusterID = "cluster4"
		want.AuthAddrInt = netx.NewAddr("authhost4int", 4104)
		want.BindAddr = netx.NewAddr("host4", 4101)
		want.BindMetricsAddr = netx.NewAddr("hostint4", 4102)
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
