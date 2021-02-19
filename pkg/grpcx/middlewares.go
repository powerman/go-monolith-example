package grpcx

import (
	"net"
	"path"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/powerman/structlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/powerman/go-monolith-example/pkg/def"
	"github.com/powerman/go-monolith-example/pkg/reflectx"
)

var (
	errUnknown  = status.Error(codes.Unknown, "unknown error")
	errInternal = status.Error(codes.Internal, "internal error")
)

// MakeUnaryServerLogger returns a new unary server interceptor that contains request logger.
func MakeUnaryServerLogger(service string, skip int) grpc.UnaryServerInterceptor {
	pkg := reflectx.CallerPkg(skip + 1)
	return func(ctx Ctx, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		log := newLogger(ctx, service, pkg, info.FullMethod)
		ctx = structlog.NewContext(ctx, log)
		return handler(ctx, req)
	}
}

// MakeStreamServerLogger returns a new stream server interceptor that contains request logger.
func MakeStreamServerLogger(service string, skip int) grpc.StreamServerInterceptor {
	pkg := reflectx.CallerPkg(skip + 1)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()
		log := newLogger(ctx, service, pkg, info.FullMethod)
		ctx = structlog.NewContext(ctx, log)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

// MakeUnaryServerRecover returns a new unary server interceptor that recover and logs panic.
func MakeUnaryServerRecover(metric def.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx Ctx, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		panicked := true
		defer func() {
			if p := recover(); panicked {
				metric.PanicsTotal.Inc()
				log := structlog.FromContext(ctx, nil)
				log.PrintErr("panic", def.LogGRPCCode, codes.Internal, "err", p,
					structlog.KeyStack, structlog.Auto)
				err = errInternal
			}
		}()
		res, err := handler(ctx, req)
		panicked = false
		return res, err
	}
}

// MakeStreamServerRecover returns a new stream server interceptor that recover and logs panic.
func MakeStreamServerRecover(metric def.Metrics) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true
		defer func() {
			if p := recover(); panicked {
				metric.PanicsTotal.Inc()
				log := structlog.FromContext(stream.Context(), nil)
				log.PrintErr("panic", "err", p, structlog.KeyStack, structlog.Auto)
				err = errInternal
			}
		}()
		err = handler(srv, stream)
		panicked = false
		return err
	}
}

// UnaryServerAccessLog returns a new unary server interceptor that logs request status.
func UnaryServerAccessLog(ctx Ctx, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	resp, err := handler(ctx, req)
	log := structlog.FromContext(ctx, nil)
	err = logHandler(log, err)
	return resp, err
}

// StreamServerAccessLog returns a new stream server interceptor that logs request status.
func StreamServerAccessLog(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	log := structlog.FromContext(stream.Context(), nil)
	log.Info("started")
	err = handler(srv, stream)
	err = logHandler(log, err)
	return err
}

// MakeUnaryClientLogger returns a new unary client interceptor that contains request logger.
func MakeUnaryClientLogger(service string, skip int) grpc.UnaryClientInterceptor {
	pkg := reflectx.CallerPkg(skip + 1)
	return func(ctx Ctx, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log := newLogger(ctx, service, pkg, method)
		ctx = structlog.NewContext(ctx, log)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// MakeStreamClientLogger returns a new stream client interceptor that contains request logger.
func MakeStreamClientLogger(service string, skip int) grpc.StreamClientInterceptor {
	pkg := reflectx.CallerPkg(skip + 1)
	return func(ctx Ctx, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log := newLogger(ctx, service, pkg, method)
		ctx = structlog.NewContext(ctx, log)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// UnaryClientAccessLog returns a new unary client interceptor that logs request status.
func UnaryClientAccessLog(ctx Ctx, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	log := structlog.FromContext(ctx, nil)
	err = logHandler(log, err)
	return err
}

// StreamClientAccessLog returns a new stream client interceptor that logs request status.
func StreamClientAccessLog(ctx Ctx, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log := structlog.FromContext(ctx, nil)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	if status.Code(err) == codes.OK {
		log.Info("started")
	} else {
		err = logHandler(log, err)
	}
	return clientStream, err
}

func newLogger(ctx Ctx, service, pkg, fullMethod string) *structlog.Logger {
	kvs := []interface{}{
		structlog.KeyApp, service,
		structlog.KeyUnit, pkg,
		def.LogFunc, path.Base(fullMethod),
		def.LogGRPCCode, "",
	}
	if p, ok := peer.FromContext(ctx); ok {
		ip, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			kvs = append(kvs, def.LogRemoteIP, ip)
		}
	}
	return structlog.New(kvs...)
}

// LogHandler logs error and hides message&details of Unknown&Internal errors.
func logHandler(log *structlog.Logger, err error) error { //nolint:funlen,gocyclo,cyclop // By design.
	s := status.Convert(err)
	code, msg := s.Code(), s.Message()
	switch code {
	case codes.OK:
		log.Info("handled", def.LogGRPCCode, code)
	case codes.Canceled:
		log.Info("handled", def.LogGRPCCode, code)
	case codes.Unknown:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
		err = errUnknown
	case codes.InvalidArgument:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.DeadlineExceeded:
		log.Warn("failed to handle", def.LogGRPCCode, code)
	case codes.NotFound:
		log.Info("handled", def.LogGRPCCode, code, "err", msg)
	case codes.AlreadyExists:
		log.Info("handled", def.LogGRPCCode, code, "err", msg)
	case codes.PermissionDenied:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.ResourceExhausted:
		log.Warn("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.FailedPrecondition:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.Aborted:
		log.Warn("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.OutOfRange:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.Unimplemented:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.Internal:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
		err = errInternal
	case codes.Unavailable:
		log.Warn("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.DataLoss:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	case codes.Unauthenticated:
		log.PrintErr("failed to handle", def.LogGRPCCode, code)
	default:
		log.PrintErr("failed to handle", def.LogGRPCCode, code, "err", msg)
	}
	return err
}
