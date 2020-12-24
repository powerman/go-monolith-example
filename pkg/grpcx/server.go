package grpcx

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/powerman/must"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
	extraUnary []grpc.UnaryServerInterceptor,
	extraStream []grpc.StreamServerInterceptor,
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
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(append([]grpc.UnaryServerInterceptor{
			serverMetrics.UnaryServerInterceptor(),
			MakeUnaryServerLogger(service, 1),
			MakeUnaryServerRecover(metric),
			UnaryServerAccessLog,
		}, extraUnary...)...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(append([]grpc.StreamServerInterceptor{
			serverMetrics.StreamServerInterceptor(),
			MakeStreamServerLogger(service, 1),
			MakeStreamServerRecover(metric),
			StreamServerAccessLog,
		}, extraStream...)...)),
	)
	reflection.Register(server)
	healthServer = health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	return server, healthServer
}

// RemoteIP returns either peer IP, or IP from X-Forwarded-For metadata
// key provided by allowed peer, or empty string if neither is available.
func RemoteIP(ctx Ctx, xffAllowedFrom func(peerIP string) bool) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		// TODO Probably it's a local call, in which case we should trust XFF..?
		return ""
	}
	remoteAddr := p.Addr.String()
	remoteIP, _, err := net.SplitHostPort(remoteAddr)
	must.NoErr(err)

	md, _ := metadata.FromIncomingContext(ctx)
	if vals := md.Get(xForwardedFor); len(vals) > 0 {
		r := &http.Request{
			RemoteAddr: remoteAddr,
			Header:     http.Header{xForwardedFor: vals},
		}
		remoteAddr = xff.GetRemoteAddrIfAllowed(r, xffAllowedFrom)
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			remoteIP = ip
		} else {
			log := structlog.FromContext(ctx, nil)
			log.Warn("failed to SplitHostPort", "xffAddr", remoteAddr)
		}
	}
	return remoteIP
}

// AccessToken returns "Bearer" AccessToken from authorization metadata, if any.
func AccessToken(ctx Ctx) string {
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get("authorization")
	const pfx = "Bearer " // OAuth require case-sensitive "Bearer", but RFC require case-insensitive https://tools.ietf.org/html/rfc7235#section-2.1
	if len(vals) > 0 && len(vals[0]) > len(pfx) && strings.EqualFold(pfx, vals[0][:len(pfx)]) {
		return vals[0][len(pfx):]
	}
	return ""
}
