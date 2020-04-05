// Package config provides a convenient way to get subsystem's configuration.
package config

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/flags"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // By design.
var shared struct {
	metricsHost   flags.NotEmptyString
	rpcHost       flags.NotEmptyString
	mysqlHost     flags.NotEmptyString
	mysqlPort     flags.Port
	natsUrls      flags.NotEmptyString
	stanClusterID flags.NotEmptyString
}

// Metrics provides host (same for everyone) and port to serve Prometheus
// metrics.
type Metrics struct {
	port flags.Port
}

// AddTo adds required flags/defaults to given cmd.
func (c *Metrics) AddTo(cmd *cobra.Command, pfx string, defPort int) {
	flags.Var(cmd, &shared.metricsHost, "metrics.host", def.MetricsHost, "serve Prometheus metrics on host")
	flags.Var(cmd, &c.port, join(pfx, "metrics.port"), defPort, "serve Prometheus metrics on port")
}

// Host returns host.
func (c Metrics) Host() string { return string(shared.metricsHost) }

// Port returns port.
func (c Metrics) Port() int { return int(c.port) }

// String returns host:port.
func (c Metrics) String() string { return fmt.Sprintf("%s:%d", c.Host(), c.Port()) }

// RPC implements (shared) rpc.host and rpc.port flags.
type RPC struct {
	port flags.Port
}

// AddTo adds required flags/defaults to given cmd.
func (c *RPC) AddTo(cmd *cobra.Command, pfx string, defPort int) {
	flags.Var(cmd, &shared.rpcHost, "rpc.host", def.RPCHost, "serve JSON-RPC 2.0 on host")
	flags.Var(cmd, &c.port, join(pfx, "rpc.port"), defPort, "serve JSON-RPC 2.0 on port")
}

// Host returns host.
func (c RPC) Host() string { return string(shared.rpcHost) }

// Port returns port.
func (c RPC) Port() int { return int(c.port) }

// String returns host:port.
func (c RPC) String() string { return fmt.Sprintf("%s:%d", c.Host(), c.Port()) }

// MySQL implements (shared) mysql.{host,port} and
// mysql.{user,pass,dbname,goose_dir} flags.
type MySQL struct {
	user     flags.NotEmptyString
	pass     string
	dbName   flags.NotEmptyString
	gooseDir flags.NotEmptyString
}

// MySQLDef provides defaults for MySQL flags.
type MySQLDef struct {
	User     string
	Pass     string
	DBName   string
	GooseDir string
}

// AddTo adds required flags/defaults to given cmd.
func (c *MySQL) AddTo(cmd *cobra.Command, pfx string, d MySQLDef) {
	cmdf := cmd.Flags()
	flags.Var(cmd, &shared.mysqlHost, "mysql.host", def.MySQLHost, "MySQL host")
	flags.Var(cmd, &shared.mysqlPort, "mysql.port", def.MySQLPort, "MySQL port")
	flags.Var(cmd, &c.user, join(pfx, "mysql.user"), d.User, "MySQL user")
	cmdf.StringVar(&c.pass, join(pfx, "mysql.pass"), d.Pass, "MySQL pass")
	flags.Var(cmd, &c.dbName, join(pfx, "mysql.dbname"), d.DBName, "MySQL database name")
	flags.Var(cmd, &c.gooseDir, join(pfx, "mysql.goose.dir"), d.GooseDir, "goose migrations dir")
}

// Config creates MySQL config.
func (c MySQL) Config() *mysql.Config {
	return def.MySQLCfg(string(shared.mysqlHost), int(shared.mysqlPort), def.MySQLAuth{
		User: string(c.user),
		Pass: c.pass,
		DB:   string(c.dbName),
	})
}

// GooseDir returns goose migrations dir.
func (c MySQL) GooseDir() string { return string(c.gooseDir) }

// NATSUrls implements (shared) nats.urls flag.
type NATSUrls struct{}

// AddTo adds required flags/defaults to given cmd.
func (c *NATSUrls) AddTo(cmd *cobra.Command) {
	flags.Var(cmd, &shared.natsUrls, "nats.urls", def.NATSUrls, "NATS urls")
}

// String returns NATS urls.
func (NATSUrls) String() string { return string(shared.natsUrls) }

// STANClusterID implements (shared) stan.cluster_id flag.
type STANClusterID struct{}

// AddTo adds required flags/defaults to given cmd.
func (c *STANClusterID) AddTo(cmd *cobra.Command) {
	flags.Var(cmd, &shared.stanClusterID, "stan.cluster_id", def.STANClusterID, "STAN cluster ID")
}

// String returns STAN cluster ID.
func (STANClusterID) String() string { return string(shared.stanClusterID) }
