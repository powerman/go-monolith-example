package def

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// MySQLAuth contains MySQL authentication details.
type MySQLAuth struct {
	User string
	Pass string
	DB   string
}

// MySQLCfg creates a new default config for MySQL.
func MySQLCfg(host string, port int, auth MySQLAuth) *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.User = auth.User
	cfg.Passwd = auth.Pass
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", host, port)
	cfg.DBName = auth.DB
	cfg.Params = map[string]string{
		"sql_mode": "'ONLY_FULL_GROUP_BY,TRADITIONAL'", // 5.7 defaults + STRICT_ALL_TABLES
	}
	cfg.Collation = "utf8mb4_unicode_ci"
	cfg.ParseTime = true
	cfg.RejectReadOnly = true
	return cfg
}

// TestMySQLCfg creates a new default config for MySQL integration tests.
func TestMySQLCfg(auth MySQLAuth) *mysql.Config {
	var connectTimeout = 3 * TestSecond

	cfg := MySQLCfg(MySQLHost, MySQLPort, auth)
	cfg.Timeout = connectTimeout
	return cfg
}
