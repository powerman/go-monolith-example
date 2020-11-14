package def

import (
	"github.com/go-sql-driver/mysql"

	"github.com/powerman/go-monolith-example/pkg/netx"
)

// MySQLConfig contains MySQL connection and authentication details.
type MySQLConfig struct {
	Addr   netx.Addr
	User   string
	Pass   string
	DBName string
}

// NewMySQLConfig creates a new default config for MySQL.
func NewMySQLConfig(cfg MySQLConfig) *mysql.Config {
	c := mysql.NewConfig()
	c.User = cfg.User
	c.Passwd = cfg.Pass
	c.Net = "tcp"
	c.Addr = cfg.Addr.String()
	c.DBName = cfg.DBName
	c.Params = map[string]string{
		"sql_mode": "'ONLY_FULL_GROUP_BY,TRADITIONAL'", // 5.7 defaults + STRICT_ALL_TABLES
	}
	c.Collation = "utf8mb4_unicode_ci"
	c.ParseTime = true
	c.RejectReadOnly = true
	return c
}
