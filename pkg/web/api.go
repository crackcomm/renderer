package web

import (
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"bitbucket.org/moovie/util/template"

	"github.com/crackcomm/renderer/pkg/renderer"
)

// New - New renderer web server API handler.
// Context should have a compiler set with `renderer.NewContext`.
//
// Default options are:
//
// 	WithTimeout(time.Second * 15),
// 	WithComponentSetter(UnmarshalFromRequest()),
//
func New(opts ...Option) xhandler.HandlerC {
	o := &options{
		reqTimeout:      time.Second * 15,
		componentSetter: UnmarshalFromRequest,
	}
	o.templateCtxSetter = defaultCtxSetter(o)

	for _, opt := range opts {
		opt(o)
	}

	var chain xhandler.Chain
	chain.UseC(xhandler.CloseHandler)
	chain.UseC(xhandler.TimeoutHandler(o.reqTimeout))
	chain.UseC(o.componentSetter)
	chain.UseC(o.templateCtxSetter)
	for _, m := range o.middlewares {
		chain.UseC(m)
	}
	chain.UseC(CompileFromCtx)
	chain.UseC(RenderFromCtx)
	if o.alwaysHTML {
		return chain.HandlerCF(WriteRenderedHTML)
	}
	return chain.HandlerCF(WriteRendered)
}

// Option - Sets web server handler options.
type Option func(*options)

type options struct {
	// reqTimeout - HTTP request timeout.
	reqTimeout time.Duration

	// componentSetter - Component setter middleware.
	// It should set a component in context using `renderer.WithComponent`.
	componentSetter Middleware

	// templateCtxSetter - Component template context setter middleware.
	// This middleware can read context from request etc.
	// It should set template context using `renderer.WithTemplateContext`.
	templateCtxSetter Middleware

	// defaultCtx - Default template context.
	// Used only when no template context setter is used.
	defaultCtx template.Context

	// alwaysHTML - Responds with html always when true.
	alwaysHTML bool

	// middlewares - List of middlewares.
	middlewares []Middleware
}

// WithComponentSetter - Sets component reader HTTP request middleware.
func WithComponentSetter(componentSetter Middleware) Option {
	return func(o *options) {
		o.componentSetter = componentSetter
	}
}

// WithTemplateCtxSetter - Sets component template context setter HTTP request middleware.
func WithTemplateCtxSetter(templateCtxSetter Middleware) Option {
	return func(o *options) {
		o.templateCtxSetter = templateCtxSetter
	}
}

// WithTimeout - Sets API server request timeout.
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.reqTimeout = t
	}
}

// WithAlwaysHTML - Responds with html only when enabled.
func WithAlwaysHTML(enable ...bool) Option {
	return func(o *options) {
		if len(enable) == 0 {
			o.alwaysHTML = true
		} else {
			o.alwaysHTML = enable[0]
		}
	}
}

// WithDefaultTemplateCtx - Sets default template context.
func WithDefaultTemplateCtx(ctx template.Context) Option {
	return func(o *options) {
		o.defaultCtx = ctx
	}
}

// WithMiddleware - Adds a middleware.
func WithMiddleware(m Middleware) Option {
	return func(o *options) {
		o.middlewares = append(o.middlewares, m)
	}
}

type breakKeyT struct{}

var breakKey = breakKeyT{}

// Break - Sets a break in context.
func Break(ctx context.Context) context.Context {
	return context.WithValue(ctx, breakKey, true)
}

// HasBreak - Returns true if break was set in given context.
func HasBreak(ctx context.Context) bool {
	b, _ := ctx.Value(breakKey).(bool)
	return b
}

func defaultCtxSetter(o *options) Middleware {
	return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		if _, ok := renderer.TemplateContext(ctx); !ok && o.defaultCtx != nil {
			ctx = renderer.WithTemplateContext(ctx, o.defaultCtx)
		}
		next.ServeHTTPC(ctx, w, r)
	})
}
