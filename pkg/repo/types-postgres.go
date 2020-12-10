package repo

import (
	"net"
	"time"

	"github.com/jackc/pgtype"
	"github.com/powerman/must"
)

// PostgresFromDuration creates new interval from duration.
func PostgresFromDuration(dur time.Duration) *pgtype.Interval {
	interval := new(pgtype.Interval)
	must.NoErr(interval.Set(dur))
	return interval
}

// PostgresInet is a wrapper around pgtype.Inet.
type PostgresInet struct {
	*pgtype.Inet
}

// PostgresFromIP creates new Inet from ip.
func PostgresFromIP(ip net.IP) *PostgresInet {
	inet := new(pgtype.Inet)
	if ip == nil || ip.IsUnspecified() {
		must.NoErr(inet.Set(nil))
	} else {
		must.NoErr(inet.Set(ip))
	}
	return &PostgresInet{Inet: inet}
}

// IP returns nil or src.IPNet.IP.
func (inet *PostgresInet) IP() net.IP {
	if inet != nil && inet.Inet != nil && inet.Inet.Status == pgtype.Present {
		return inet.IPNet.IP
	}
	return nil
}

// Scan implements the database/sql Scanner interface.
func (inet *PostgresInet) Scan(src interface{}) error {
	if inet.Inet == nil {
		inet.Inet = new(pgtype.Inet)
	}
	return inet.Inet.Scan(src)
}
