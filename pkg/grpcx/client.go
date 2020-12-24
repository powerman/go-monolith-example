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
	"google.golang.org/grpc/metadata"
)

// Dial creates a gRPC client connection to the given target.
func Dial(ctx Ctx, addr, service string, metrics *grpc_prometheus.ClientMetrics, ca *x509.CertPool) (*grpc.ClientConn, error) {
	opts := append(DialOptions(ca),
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
	return grpc.DialContext(ctx, addr, opts...)
}

// DialOptions returns default connection options without interceptors.
func DialOptions(ca *x509.CertPool) []grpc.DialOption {
	const serviceConfigHealthCheck = `{
		"loadBalancingPolicy": "round_robin",
		"healthCheckConfig": {
			"serviceName": ""
		},
		"methodConfig": [{
			"name": [{"service":""}],
			"waitForReady": true
		}]
	}`
	return []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(ca, "")),
		grpc.WithDefaultServiceConfig(serviceConfigHealthCheck),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveTime,
			Timeout:             keepaliveTimeout,
			PermitWithoutStream: true,
		}),
	}
}

// AppendXFF returns a new context with the provided X-Forwarded-For value
// merged with any existing metadata in the outgoing context.
func AppendXFF(ctx Ctx, xff string) Ctx {
	return metadata.AppendToOutgoingContext(ctx, xForwardedFor, xff)
}

// AccessTokenCreds returns a CallOption that sets
// credentials.PerRPCCredentials using OAuth2 "Bearer" AccessToken.
func AccessTokenCreds(accessToken string) grpc.CallOption {
	var creds credentials.PerRPCCredentials
	if accessToken != "" {
		creds = oauth.NewOauthAccess(&oauth2.Token{AccessToken: accessToken})
	}
	return grpc.PerRPCCredentials(creds)
}
