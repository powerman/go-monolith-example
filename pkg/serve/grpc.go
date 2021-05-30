package serve

import (
	"net"

	"github.com/powerman/structlog"
	"google.golang.org/grpc"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/netx"
)

// GRPC starts gRPC server on addr, logged as service.
// It runs until failed or ctx.Done.
func GRPC(ctx Ctx, addr netx.Addr, srv *grpc.Server, service string) (err error) {
	log := structlog.FromContext(ctx, nil).New(def.LogServer, service)

	ln, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	log.Info("serve", def.LogHost, addr.Host(), def.LogPort, addr.Port())
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()

	select {
	case err = <-errc:
	case <-ctx.Done():
		srv.GracefulStop() // It will not interrupt streaming.
	}
	if err != nil {
		return log.Err("failed to serve", "err", err)
	}
	log.Info("shutdown")
	return nil
}
