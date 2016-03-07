package renderer

import (
	"errors"

	"golang.org/x/net/context"

	"bitbucket.org/moovie/util/template"
)

type compilerCtxKeyT struct{}

var compilerCtxKey = compilerCtxKeyT{}

// CompilerCtx - Creates a new context with compiler.
func CompilerCtx(ctx context.Context, c Compiler) context.Context {
	return context.WithValue(ctx, compilerCtxKey, c)
}

// CompilerFromCtx - Retrieves compiler from context.
func CompilerFromCtx(ctx context.Context) (c Compiler, ok bool) {
	c, ok = ctx.Value(compilerCtxKey).(Compiler)
	return
}

type compiledCtxKeyT struct{}

var compiledCtxKey = compiledCtxKeyT{}

// CompiledCtx - Creates a new context with compiled.
func CompiledCtx(ctx context.Context, c *Compiled) context.Context {
	return context.WithValue(ctx, compiledCtxKey, c)
}

// CompiledFromCtx - Retrieves compiled from context.
func CompiledFromCtx(ctx context.Context) (c *Compiled, ok bool) {
	c, ok = ctx.Value(compiledCtxKey).(*Compiled)
	return
}

type componentCtxKeyT struct{}

var componentCtxKey = componentCtxKeyT{}

// ComponentCtx - Creates a new context with component.
func ComponentCtx(ctx context.Context, c *Component) context.Context {
	return context.WithValue(ctx, componentCtxKey, c)
}

// ComponentFromCtx - Retrieves component from context.
func ComponentFromCtx(ctx context.Context) (c *Component, ok bool) {
	c, ok = ctx.Value(componentCtxKey).(*Component)
	return
}

type renderedCtxKeyT struct{}

var renderedCtxKey = renderedCtxKeyT{}

// RenderedCtx - Creates a new context with rendered component.
func RenderedCtx(ctx context.Context, c *Rendered) context.Context {
	return context.WithValue(ctx, renderedCtxKey, c)
}

// RenderedFromCtx - Retrieves rendered component from context.
func RenderedFromCtx(ctx context.Context) (c *Rendered, ok bool) {
	c, ok = ctx.Value(renderedCtxKey).(*Rendered)
	return
}

type templateCtxKeyT struct{}

var templateCtxKey = templateCtxKeyT{}

// WithTemplateCtx - Creates a new context with template context set.
// Template context can be retrieved using `renderer.TemplateCtx(context.Context)`.
func WithTemplateCtx(ctx context.Context, t template.Context) context.Context {
	return context.WithValue(ctx, templateCtxKey, t)
}

// TemplateCtx - Retrieves `template.Context` from `context.Context`.
func TemplateCtx(ctx context.Context) (t template.Context, ok bool) {
	t, ok = ctx.Value(templateCtxKey).(template.Context)
	return
}

// SetTemplateCtx - Sets template context value in `context.Context`.
func SetTemplateCtx(ctx context.Context, key string, v interface{}) context.Context {
	t, ok := TemplateCtx(ctx)
	if ok {
		t[key] = v
		return ctx
	}
	return WithTemplateCtx(ctx, template.Context{key: v})
}

// CompileFromCtx - Compiles component from context
// using compiler and storage from context.
func CompileFromCtx(ctx context.Context) (compiled *Compiled, err error) {
	c, ok := ComponentFromCtx(ctx)
	if !ok {
		err = errors.New("no component set")
		return
	}
	compiler, ok := CompilerFromCtx(ctx)
	if !ok {
		err = errors.New("compiler not found")
		return
	}
	return compiler.CompileFromStorage(c)
}
