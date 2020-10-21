// Package natsx implements connections to NATS and NATS Streaming (STAN).
package natsx

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/powerman/structlog"
)

const (
	drainTimeout  = 3 * time.Second // Should be less than main.shutdownDelay.
	maxReconnects = 5
	pingInterval  = time.Second // Default 2 min isn't useful because TCP keepalive is faster.
)

// NATSConn adds connection monitoring to nats.Conn.
type NATSConn struct {
	*nats.Conn
	closed chan struct{}
	log    *structlog.Logger
}

// ConnectNATS adds ctx support and reasonable defaults to nats.Connect.
func ConnectNATS(ctx context.Context, urls, name string) (*NATSConn, error) {
	c := &NATSConn{
		closed: make(chan struct{}),
		log:    structlog.FromContext(ctx, nil),
	}
	var err error
	for !(c.Conn != nil && err == nil) {
		errc := make(chan error)
		go func() {
			err := c.connect(urls, name)
			select {
			case errc <- err:
			case <-ctx.Done():
				if c.Conn != nil {
					c.Conn.Close()
				}
			}
		}()
		select {
		case err = <-errc:
		case <-ctx.Done():
			if err == nil {
				err = ctx.Err()
			}
			return nil, err
		}
	}
	c.log.Info("NATS connected", "url", c.ConnectedUrl())
	return c, nil
}

func (c *NATSConn) connect(urls, name string) (err error) {
	c.Conn, err = nats.Connect(urls,
		nats.Name(name),
		nats.MaxReconnects(maxReconnects),
		nats.DrainTimeout(drainTimeout),
		nats.PingInterval(pingInterval),
		nats.NoCallbacksAfterClientClose(),
		nats.ClosedHandler(func(_ *nats.Conn) {
			close(c.closed)
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err == nil {
				c.log.Info("NATS disconnected")
			} else {
				c.log.Warn("NATS disconnected", "err", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			c.log.Info("NATS reconnected", "url", nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(_ *nats.Conn, sub *nats.Subscription, err error) {
			if sub == nil {
				c.log.Warn("NATS connection failed", "err", err)
			} else {
				c.log.Warn("NATS connection failed", "subject", sub.Subject, "err", err)
			}
		}),
	)
	return err
}

// Monitor waits until ctx.Done or failure reconnecting NATS.
func (c *NATSConn) Monitor(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case <-c.closed:
		return c.log.Err("NATS connection closed: unable to connect")
	}
}

// STANConn adds connection monitoring to stan.Conn.
type STANConn struct {
	stan.Conn
	closed chan error
	log    *structlog.Logger
}

// ConnectSTAN adds ctx support and reasonable defaults to stan.Connect.
func ConnectSTAN(ctx context.Context, clusterID, clientID string, nc *NATSConn) (*STANConn, error) {
	clientID = regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(clientID, "-")

	c := &STANConn{
		closed: make(chan error, 1),
		log:    structlog.FromContext(ctx, nil),
	}
	err := c.connect(ctx, clusterID, clientID, nc)
	if err != nil {
		return nil, err
	}
	c.log.Info("STAN connected", "clusterID", clusterID, "clientID", clientID)
	return c, nil
}

func (c *STANConn) connect(ctx context.Context, clusterID, clientID string, nc *NATSConn) (err error) {
	connectWait := stan.DefaultConnectWait
	if deadline, ok := ctx.Deadline(); ok {
		connectWait = time.Until(deadline)
	}
	c.Conn, err = stan.Connect(clusterID, clientID,
		stan.NatsConn(nc.Conn),
		stan.ConnectWait(connectWait),
		stan.Pings(int(pingInterval.Seconds()), stan.DefaultPingMaxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			c.closed <- fmt.Errorf("STAN connection closed: %w", err)
			close(c.closed)
		}),
	)
	return err
}

// Monitor waits until ctx.Done or closed STAN connection.
func (c *STANConn) Monitor(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case err := <-c.closed:
		return c.log.Err(err)
	}
}
