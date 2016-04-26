package compiler

import (
	"errors"

	"tower.pro/renderer/components"

	"golang.org/x/net/context"
)

type compilerCtxKey struct{}

var compilerKey = compilerCtxKey{}

// NewContext - Creates a new context with compiler.
func NewContext(ctx context.Context, c *Compiler) context.Context {
	return context.WithValue(ctx, compilerKey, c)
}

// FromContext - Retrieves compiler from context.
func FromContext(ctx context.Context) (c *Compiler, ok bool) {
	c, ok = ctx.Value(compilerKey).(*Compiler)
	return
}

// Compile - Compiles component from context using compiler from context.
func Compile(ctx context.Context) (compiled *components.Compiled, err error) {
	c, ok := components.FromContext(ctx)
	if !ok {
		err = errors.New("no component set")
		return
	}
	compiler, ok := FromContext(ctx)
	if !ok {
		err = errors.New("compiler not found")
		return
	}
	return compiler.CompileFromStorage(c)
}
