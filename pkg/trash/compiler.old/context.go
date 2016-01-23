package compiler

import "golang.org/x/net/context"

type compilerCtxKey struct{}

// NewContext - Creates new context with compiler.
func NewContext(ctx context.Context, c Compiler) context.Context {
	return context.WithValue(ctx, compilerCtxKey{}, c)
}

// FromContext - Retrieves compiler from context.
func FromContext(ctx context.Context) (c Compiler, ok bool) {
	c, ok = ctx.Value(compilerCtxKey{}).(Compiler)
	return
}
