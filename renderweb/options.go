package renderweb

import (
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/template"

	"github.com/crackcomm/renderer/components"
	"github.com/crackcomm/renderer/middlewares"
)

// Option - Sets web server handler options.
type Option func(*options)

type options struct {
	alwaysHTML bool
	reqTimeout time.Duration
	defaultCtx template.Context

	middlewares       []middlewares.Handler
	componentSetter   middlewares.Handler
	templateCtxSetter middlewares.Handler
}

// WithComponentSetter - Sets component reader HTTP request middleware.
func WithComponentSetter(componentSetter middlewares.Handler) Option {
	return func(o *options) {
		o.componentSetter = componentSetter
	}
}

// WithTemplateContextSetter - Sets component template context setter HTTP request middleware.
func WithTemplateContextSetter(templateCtxSetter middlewares.Handler) Option {
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

// WithDefaultTemplateContext - Sets default template context.
func WithDefaultTemplateContext(ctx template.Context) Option {
	return func(o *options) {
		o.defaultCtx = ctx
	}
}

// WithMiddleware - Adds a middleware.
func WithMiddleware(m middlewares.Handler) Option {
	return func(o *options) {
		o.middlewares = append(o.middlewares, m)
	}
}

func defaultCtxSetter(o *options) middlewares.Handler {
	return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		if _, ok := components.TemplateContext(ctx); !ok && o.defaultCtx != nil {
			ctx = components.NewTemplateContext(ctx, o.defaultCtx.Clone())
		}
		next.ServeHTTPC(ctx, w, r)
	})
}
