package apix

import (
	"context"
	"path"
	"strings"

	"github.com/powerman/structlog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/pkg/def"
)

// GRPCNewContext returns a new context.Context that carries values describing
// this request without any deadline, plus result of authn.Authenticate.
func GRPCNewContext(ctx Ctx, fullMethod string, authn Authn) (_ Ctx, auth dom.Auth, err error) {
	if p, ok := peer.FromContext(ctx); ok {
		ctx = context.WithValue(ctx, contextKeyRemote, p.Addr.String())
	}

	ctx = context.WithValue(ctx, contextKeyMethodName, path.Base(fullMethod))

	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get("authorization")
	if len(vals) > 0 && strings.HasPrefix(vals[0], "Bearer ") {
		accessToken := AccessToken(strings.TrimPrefix(vals[0], "Bearer "))
		ctx = context.WithValue(ctx, contextKeyAccessToken, accessToken)
		auth, err = authn.Authenticate(ctx, accessToken)
		if err == nil {
			ctx = context.WithValue(ctx, contextKeyAuth, auth)
			structlog.FromContext(ctx, nil).SetDefaultKeyvals(def.LogUserName, auth.UserName)
		}
	}

	return ctx, auth, err
}
