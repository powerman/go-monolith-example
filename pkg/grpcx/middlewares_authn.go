package grpcx

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// AuthnFunc provides a way to check authentication using interceptor.
// FullMethod is the full RPC method string, i.e., /package.service/method.
// It usually either returns error with codes.Unauthenticated or Ctx with
// extra value describing current authentication for use in handlers.
type AuthnFunc func(_ Ctx, fullMethod string) (Ctx, error)

// MakeUnaryServerAuthn returns a new unary server interceptor that checks authentication.
func MakeUnaryServerAuthn(authn AuthnFunc) grpc.UnaryServerInterceptor {
	return func(ctx Ctx, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		ctx, err = authn(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// MakeStreamServerAuthn returns a new stream server interceptor that checks authentication.
func MakeStreamServerAuthn(authn AuthnFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()
		ctx, err = authn(ctx, info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}
