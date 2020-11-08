package grpcx

import (
	"crypto/tls"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/powerman/go-monolith-example/pkg/def"
)

// NewServer creates and returns a gRPC server which:
//   - has configured TLS,
//   - has configured keep-alive,
//   - setup interceptor to provide prometheus metrics,
//   - setup interceptor to store request-scooped logger inside context,
//   - setup interceptor to recover from panics,
//   - setup interceptor to log method access/result,
//   - setup interceptor to check authentication using given authn,
//   - has reflection service registered,
//   - has health service registered,
//   - has not started to accept requests yet.
// It also returns health server which may be used to control status
// returned by health service (set to SERVING by default).
func NewServer(
	service string,
	metric def.Metrics,
	serverMetrics *grpc_prometheus.ServerMetrics,
	cert *tls.Certificate,
	authn AuthnFunc,
) (server *grpc.Server, healthServer *health.Server) {
	server = grpc.NewServer(
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
	reflection.Register(server)
	healthServer = health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	return server, healthServer
}
