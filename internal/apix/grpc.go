package apix

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/powerman/structlog"
	"github.com/sebest/xff"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/def"
)

const (
	grpcGatewayIP = "127.0.0.1"
	xForwardedFor = "X-Forwarded-For"
)

func isGRPCGateway(sip string) bool { return sip == grpcGatewayIP }

// GRPCNewContext returns a new context.Context that carries values describing
// this request without any deadline, plus result of authn.Authenticate.
func GRPCNewContext(ctx Ctx, fullMethod string, authn Authn) (_ Ctx, auth dom.Auth, err error) {
	md, _ := metadata.FromIncomingContext(ctx)

	if p, ok := peer.FromContext(ctx); ok {
		remote := p.Addr.String()
		if vals := md.Get(xForwardedFor); len(vals) > 0 {
			r := &http.Request{
				RemoteAddr: p.Addr.String(),
				Header:     http.Header{xForwardedFor: vals},
			}
			remote = xff.GetRemoteAddrIfAllowed(r, isGRPCGateway)
			structlog.FromContext(ctx, nil).SetDefaultKeyvals(def.LogRemote, remote)
		}
		ctx = context.WithValue(ctx, contextKeyRemote, remote)
	}

	ctx = context.WithValue(ctx, contextKeyMethodName, path.Base(fullMethod))

	vals := md.Get("authorization")
	const pfx = "Bearer " // OAuth require case-sensitive "Bearer", but RFC require case-insensitive https://tools.ietf.org/html/rfc7235#section-2.1
	if len(vals) > 0 && len(vals[0]) > len(pfx) && strings.EqualFold(pfx, vals[0][:len(pfx)]) {
		accessToken := AccessToken(vals[0][len(pfx):])
		ctx = context.WithValue(ctx, contextKeyAccessToken, accessToken)
		auth, err = authn.Authenticate(ctx, accessToken)
		if err == nil {
			ctx = context.WithValue(ctx, contextKeyAuth, auth)
			structlog.FromContext(ctx, nil).SetDefaultKeyvals(def.LogUserName, auth.UserName)
		}
	}

	return ctx, auth, err
}
