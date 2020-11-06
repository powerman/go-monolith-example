package grpcx

import (
	"crypto/tls"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/powerman/go-monolith-example/pkg/def"
)

// NewServer creates a gRPC server which has no service registered and has
// not started to accept requests yet.
func NewServer(
	service string,
	metric def.Metrics,
	serverMetrics *grpc_prometheus.ServerMetrics,
	cert *tls.Certificate,
	authn AuthnFunc,
) *grpc.Server {
	return grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(cert)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    keepaliveTime,
			Timeout: keepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepaliveMinTime,
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			serverMetrics.UnaryServerInterceptor(),
			MakeUnaryServerLogger(service, 1),
			MakeUnaryServerRecover(metric),
			UnaryServerAccessLog,
			MakeUnaryServerAuthn(authn),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			serverMetrics.StreamServerInterceptor(),
			MakeStreamServerLogger(service, 1),
			MakeStreamServerRecover(metric),
			StreamServerAccessLog,
			MakeStreamServerAuthn(authn),
		)),
	)
}
