package def

import (
	"context"
	"database/sql"
	"time"

	"github.com/powerman/pqx"
)

// Default timeouts for PostgreSQL.
const (
	PostgresDefaultStatementTimeout                = 3 * time.Second
	PostgresDefaultLockTimeout                     = 3 * time.Second
	PostgresDefaultIdleInTransactionSessionTimeout = 10 * time.Second
)

// PostgresConfig described connection parameters for github.com/lib/pq.
type PostgresConfig struct {
	pqx.Config
}

// NewPostgresConfig creates a new default config for PostgreSQL.
func NewPostgresConfig(cfg pqx.Config) *PostgresConfig {
	c := (&PostgresConfig{Config: cfg}).Clone()
	// Enforce SSL.
	if c.SSLMode != pqx.SSLVerifyFull {
		c.SSLMode = pqx.SSLVerifyCA
	}
	// Extra protection in case secure schema usage pattern isn't used on this server
	// https://www.postgresql.org/docs/11/ddl-schemas.html#DDL-SCHEMAS-PATTERNS.
	c.SearchPath = `"$user"`
	// In modern PostgreSQL serializable is fast enough, use it by default.
	if c.DefaultTransactionIsolation == sql.LevelDefault {
		c.DefaultTransactionIsolation = sql.LevelSerializable
	}
	// Sane timeout defaults:
	if c.StatementTimeout == 0 {
		c.StatementTimeout = PostgresDefaultStatementTimeout
	}
	if c.LockTimeout == 0 {
		c.LockTimeout = PostgresDefaultLockTimeout
	}
	if c.IdleInTransactionSessionTimeout == 0 {
		c.IdleInTransactionSessionTimeout = PostgresDefaultIdleInTransactionSessionTimeout
	}
	return c
}

// Clone returns a deep copy.
func (c *PostgresConfig) Clone() *PostgresConfig {
	clone := *c
	clone.Other = make(map[string]string, len(c.Other))
	for k, v := range c.Other {
		clone.Other[k] = v
	}
	return &clone
}

// UpdateConnectTimeout updates c accordingly to ctx.Deadline if
// c.ConnectTimeout isn't set or larger than ctx.Deadline.
func (c *PostgresConfig) UpdateConnectTimeout(ctx context.Context) error {
	if deadline, ok := ctx.Deadline(); ok {
		if c.ConnectTimeout == 0 || time.Until(deadline) < c.ConnectTimeout {
			c.ConnectTimeout = time.Until(deadline)
			if c.ConnectTimeout <= 0 {
				return context.DeadlineExceeded
			}
		}
	}
	return nil
}
