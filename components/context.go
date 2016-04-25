package components

import (
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/template"
)

type contextKey struct {
	name string
}

var (
	compiledCtxKey  = &contextKey{"renderer.compiled"}
	componentCtxKey = &contextKey{"renderer.component"}
	renderedCtxKey  = &contextKey{"renderer.rendered"}
	templateCtxKey  = &contextKey{"renderer.template"}
)

// NewContext - Creates a new context with component.
func NewContext(ctx context.Context, c *Component) context.Context {
	return context.WithValue(ctx, componentCtxKey, c)
}

// FromContext - Gets component from context.
func FromContext(ctx context.Context) (c *Component, ok bool) {
	c, ok = ctx.Value(componentCtxKey).(*Component)
	return
}

// NewCompiledContext - Creates a new context with compiled.
func NewCompiledContext(ctx context.Context, c *Compiled) context.Context {
	return context.WithValue(ctx, compiledCtxKey, c)
}

// NewRenderedContext - Creates a new context with rendered component.
func NewRenderedContext(ctx context.Context, c *Rendered) context.Context {
	return context.WithValue(ctx, renderedCtxKey, c)
}

// NewTemplateContext - Creates a new context with template context set.
// Template context can be retrieved using `renderer.TemplateContext(context.Context)`.
func NewTemplateContext(ctx context.Context, t template.Context) context.Context {
	return &templateContext{
		Context:     ctx,
		templateCtx: t,
	}
}

type templateContext struct {
	context.Context
	templateCtx template.Context
}

func (ctx *templateContext) Value(key interface{}) interface{} {
	if key == templateCtxKey {
		return ctx.templateCtx
	}
	keystr, ok := key.(string)
	if !ok {
		return ctx.Context.Value(key)
	}
	keysplit := strings.Split(keystr, ".")
	if len(keysplit) < 2 || keysplit[0] != "template" {
		return ctx.Context.Value(key)
	}
	value := getDeepValue(ctx.templateCtx, keysplit[1:])
	if value != nil {
		return value
	}
	return ctx.Context.Value(key)
}

func getDeepValue(v interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return v
	}
	key := keys[0]
	rest := keys[1:]
	switch t := v.(type) {
	case template.Context:
		return getDeepValue(t[key], rest)
	case map[string]interface{}:
		return getDeepValue(t[key], rest)
	case map[interface{}]interface{}:
		return getDeepValue(t[key], rest)
	case []interface{}:
		n, err := strconv.Atoi(key)
		if err != nil {
			return nil
		}
		if len(t) <= n {
			return nil
		}
		return t[n]
	}
	return nil
}

// WithTemplateKey - Sets template context key-value pair in `context.Context`.
func WithTemplateKey(ctx context.Context, key string, v interface{}) context.Context {
	t, ok := TemplateContext(ctx)
	if ok {
		t[key] = v
		return ctx
	}
	return NewTemplateContext(ctx, template.Context{key: v})
}

// TemplateContext - Retrieves `template.Context` from `context.Context`.
func TemplateContext(ctx context.Context) (t template.Context, ok bool) {
	t, ok = ctx.Value(templateCtxKey).(template.Context)
	return
}

// TemplateValue - Retrieves value from `template.Context` from `context.Context`.
func TemplateValue(ctx context.Context, key string) (v interface{}, ok bool) {
	t, ok := TemplateContext(ctx)
	if !ok {
		return
	}
	v, ok = t[key]
	return
}

// CompiledFromContext - Retrieves compiled from context.
func CompiledFromContext(ctx context.Context) (c *Compiled, ok bool) {
	c, ok = ctx.Value(compiledCtxKey).(*Compiled)
	return
}

// RenderedFromContext - Retrieves rendered component from context.
func RenderedFromContext(ctx context.Context) (c *Rendered, ok bool) {
	c, ok = ctx.Value(renderedCtxKey).(*Rendered)
	return
}
