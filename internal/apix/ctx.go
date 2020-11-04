package apix

import "context"

type contextKey int

const (
	_ contextKey = iota
	contextKeyRemote
	contextKeyMethodName
)

// FromContext returns values describing request stored in ctx, if any.
func FromContext(ctx context.Context) (remote, methodName string) {
	remote, _ = ctx.Value(contextKeyRemote).(string)
	methodName, _ = ctx.Value(contextKeyMethodName).(string)
	return remote, methodName
}
