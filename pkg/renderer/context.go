package renderer

import "golang.org/x/net/context"

type compilerCtxKey struct{}

// CompilerCtx - Creates a new context with compiler.
func CompilerCtx(ctx context.Context, c Compiler) context.Context {
	return context.WithValue(ctx, compilerCtxKey{}, c)
}

// CompilerFromCtx - Retrieves compiler from context.
func CompilerFromCtx(ctx context.Context) (c Compiler, ok bool) {
	c, ok = ctx.Value(compilerCtxKey{}).(Compiler)
	return
}

type compiledCtxKey struct{}

// CompiledCtx - Creates a new context with compiled.
func CompiledCtx(ctx context.Context, c *Compiled) context.Context {
	return context.WithValue(ctx, compiledCtxKey{}, c)
}

// CompiledFromCtx - Retrieves compiled from context.
func CompiledFromCtx(ctx context.Context) (c *Compiled, ok bool) {
	c, ok = ctx.Value(compiledCtxKey{}).(*Compiled)
	return
}

type componentCtxKey struct{}

// ComponentCtx - Creates a new context with component.
func ComponentCtx(ctx context.Context, c *Component) context.Context {
	return context.WithValue(ctx, componentCtxKey{}, c)
}

// ComponentFromCtx - Retrieves component from context.
func ComponentFromCtx(ctx context.Context) (c *Component, ok bool) {
	c, ok = ctx.Value(componentCtxKey{}).(*Component)
	return
}

type renderedCtxKey struct{}

// RenderedCtx - Creates a new context with rendered component.
func RenderedCtx(ctx context.Context, c *Rendered) context.Context {
	return context.WithValue(ctx, renderedCtxKey{}, c)
}

// RenderedFromCtx - Retrieves rendered component from context.
func RenderedFromCtx(ctx context.Context) (c *Rendered, ok bool) {
	c, ok = ctx.Value(renderedCtxKey{}).(*Rendered)
	return
}
