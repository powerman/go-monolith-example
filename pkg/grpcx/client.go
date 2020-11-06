package grpcx

import (
	"crypto/x509"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/keepalive"
)

// Dial creates a gRPC client connection to the given target.
func Dial(ctx Ctx, addr, service string, metrics *grpc_prometheus.ClientMetrics, ca *x509.CertPool) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(ca, "")),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveTime,
			Timeout:             keepaliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			metrics.UnaryClientInterceptor(),
			MakeUnaryClientLogger(service, 1),
			UnaryClientAccessLog,
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			metrics.StreamClientInterceptor(),
			MakeStreamClientLogger(service, 1),
			StreamClientAccessLog,
		)),
	)
}

// Token returns option with "Bearer" token.
func Token(token string) grpc.CallOption {
	perRPC := oauth.NewOauthAccess(&oauth2.Token{AccessToken: token})
	if token == "" {
		perRPC = nil
	}
	return grpc.PerRPCCredentials(perRPC)
}
