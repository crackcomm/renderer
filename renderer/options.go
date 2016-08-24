package renderer

import (
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"tower.pro/renderer/components"
	"tower.pro/renderer/middlewares"
	"tower.pro/renderer/template"
)

type webOptions struct {
	tracing    bool
	alwaysHTML bool
	reqTimeout time.Duration
	defaultCtx template.Context

	middlewares       []middlewares.Handler
	componentSetter   middlewares.Handler
	templateCtxSetter middlewares.Handler
}

// Option - Sets web server handler options.
type Option func(*webOptions)

func constructOpts(opts ...Option) *webOptions {
	o := &webOptions{
		reqTimeout:      time.Second * 15,
		componentSetter: UnmarshalFromRequest,
	}
	o.templateCtxSetter = defaultCtxSetter(o)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithTracing - Enables tracing. Uses first parameter if any.
func WithTracing(enable ...bool) Option {
	return func(o *webOptions) {
		if len(enable) == 0 {
			o.tracing = true
		} else {
			o.tracing = enable[0]
		}
	}
}

// WithComponentSetter - Sets component reader HTTP request middleware.
func WithComponentSetter(componentSetter middlewares.Handler) Option {
	return func(o *webOptions) {
		o.componentSetter = componentSetter
	}
}

// WithTemplateContextSetter - Sets component template context setter HTTP request middleware.
func WithTemplateContextSetter(templateCtxSetter middlewares.Handler) Option {
	return func(o *webOptions) {
		o.templateCtxSetter = templateCtxSetter
	}
}

// WithTimeout - Sets API server request timeout.
func WithTimeout(t time.Duration) Option {
	return func(o *webOptions) {
		o.reqTimeout = t
	}
}

// WithAlwaysHTML - Responds with html only when enabled. Uses first parameter if any.
func WithAlwaysHTML(enable ...bool) Option {
	return func(o *webOptions) {
		if len(enable) == 0 {
			o.alwaysHTML = true
		} else {
			o.alwaysHTML = enable[0]
		}
	}
}

// WithDefaultTemplateContext - Sets default template context.
// Template context is cloned on every request because it's used as a base.
func WithDefaultTemplateContext(ctx template.Context) Option {
	return func(o *webOptions) {
		o.defaultCtx = ctx
	}
}

// WithMiddleware - Adds a middleware.
func WithMiddleware(m middlewares.Handler) Option {
	return func(o *webOptions) {
		o.middlewares = append(o.middlewares, m)
	}
}

func defaultCtxSetter(o *webOptions) middlewares.Handler {
	return middlewares.ToHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		ctx = components.NewTemplateContext(ctx, o.defaultCtx.Clone())
		next.ServeHTTPC(ctx, w, r)
	})
}
