package renderer

import (
	"errors"

	"golang.org/x/net/context"

	"bitbucket.org/moovie/util/template"
)

type compilerCtxKeyT struct{}

var compilerCtxKey = compilerCtxKeyT{}

// WithCompiler - Creates a new context with compiler.
func WithCompiler(ctx context.Context, c Compiler) context.Context {
	return context.WithValue(ctx, compilerCtxKey, c)
}

// ContextCompiler - Retrieves compiler from context.
func ContextCompiler(ctx context.Context) (c Compiler, ok bool) {
	c, ok = ctx.Value(compilerCtxKey).(Compiler)
	return
}

type compiledCtxKeyT struct{}

var compiledCtxKey = compiledCtxKeyT{}

// WithCompiled - Creates a new context with compiled.
func WithCompiled(ctx context.Context, c *Compiled) context.Context {
	return context.WithValue(ctx, compiledCtxKey, c)
}

// ContextCompiled - Retrieves compiled from context.
func ContextCompiled(ctx context.Context) (c *Compiled, ok bool) {
	c, ok = ctx.Value(compiledCtxKey).(*Compiled)
	return
}

type componentCtxKeyT struct{}

var componentCtxKey = componentCtxKeyT{}

// WithComponent - Creates a new context with component.
func WithComponent(ctx context.Context, c *Component) context.Context {
	return context.WithValue(ctx, componentCtxKey, c)
}

// ContextComponent - Retrieves component from context.
func ContextComponent(ctx context.Context) (c *Component, ok bool) {
	c, ok = ctx.Value(componentCtxKey).(*Component)
	return
}

type renderedCtxKeyT struct{}

var renderedCtxKey = renderedCtxKeyT{}

// WithRendered - Creates a new context with rendered component.
func WithRendered(ctx context.Context, c *Rendered) context.Context {
	return context.WithValue(ctx, renderedCtxKey, c)
}

// ContextRendered - Retrieves rendered component from context.
func ContextRendered(ctx context.Context) (c *Rendered, ok bool) {
	c, ok = ctx.Value(renderedCtxKey).(*Rendered)
	return
}

type templateCtxKeyT struct{}

var templateCtxKey = templateCtxKeyT{}

// WithTemplateContext - Creates a new context with template context set.
// Template context can be retrieved using `renderer.TemplateContext(context.Context)`.
func WithTemplateContext(ctx context.Context, t template.Context) context.Context {
	return context.WithValue(ctx, templateCtxKey, t)
}

// TemplateContext - Retrieves `template.Context` from `context.Context`.
func TemplateContext(ctx context.Context) (t template.Context, ok bool) {
	t, ok = ctx.Value(templateCtxKey).(template.Context)
	return
}

// WithTemplateKey - Sets template context key-value pair in `context.Context`.
func WithTemplateKey(ctx context.Context, key string, v interface{}) context.Context {
	t, ok := TemplateContext(ctx)
	if ok {
		t[key] = v
		return ctx
	}
	return WithTemplateContext(ctx, template.Context{key: v})
}

// CompileContext - Compiles component from context
// using compiler and storage from context.
func CompileContext(ctx context.Context) (compiled *Compiled, err error) {
	c, ok := ContextComponent(ctx)
	if !ok {
		err = errors.New("no component set")
		return
	}
	compiler, ok := ContextCompiler(ctx)
	if !ok {
		err = errors.New("compiler not found")
		return
	}
	return compiler.CompileFromStorage(c)
}
