package middlewares

import "golang.org/x/net/context"

type breakCtxKey struct{}

var breakKey = breakCtxKey{}

// WithBreak - Returns new context with a break for consideration for next middlewares
func WithBreak(ctx context.Context) context.Context {
	return context.WithValue(ctx, breakKey, true)
}

// HasBreak - Returns true if break was set in given context.
func HasBreak(ctx context.Context) bool {
	b, _ := ctx.Value(breakKey).(bool)
	return b
}
