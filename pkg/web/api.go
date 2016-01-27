package web

import (
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// NewAPI - New renderer web server API handler.
// Context should have a compiler set with `renderer.NewContext`.
//
// Default options are:
//
// 	WithContext(context.Background()),
// 	WithTimeout(time.Second * 15),
// 	WithComponentSetter(UnmarshalFromRequest()),
//
func NewAPI(opts ...Option) http.Handler {
	o := &options{
		reqTimeout: time.Second * 15,
		ctx:        context.Background(),
	}

	for _, opt := range opts {
		opt(o)
	}

	// Create a Chain of middlewares
	var chain xhandler.Chain

	// Add close notifier handler so context is cancelled when the client closes
	// the connection
	chain.UseC(xhandler.CloseHandler)

	// Add timeout handler
	chain.UseC(xhandler.TimeoutHandler(o.reqTimeout))

	// Construct API from middlewares
	if o.componentSetter == nil {
		o.componentSetter = UnmarshalFromRequest()
	}
	chain.UseC(o.componentSetter)
	chain.UseC(CompileFromCtx)
	chain.UseC(RenderFromCtx)

	return chain.HandlerCtx(o.ctx, xhandler.HandlerFuncC(WriteRendered))
}

// Option - Sets web server handler options.
type Option func(*options)

type options struct {
	reqTimeout time.Duration
	ctx        context.Context

	// componentSetter - Component setter middleware.
	// It should set a component in context using `renderer.ComponentCtx`.
	componentSetter Middleware

	// templateCtxSetter - Component template context setter middleware.
	// This middleware can read context from request etc.
	// It should set template context using `renderer.WithTemplateCtx`.
	templateCtxSetter Middleware
}

// WithComponentSetter - Sets component reader HTTP request middleware.
func WithComponentSetter(componentSetter Middleware) Option {
	return func(o *options) {
		o.componentSetter = componentSetter
	}
}

// WithCtxSetter - Sets component template context setter HTTP request middleware.
func WithCtxSetter(templateCtxSetter Middleware) Option {
	return func(o *options) {
		o.templateCtxSetter = templateCtxSetter
	}
}

// WithContext - Sets API server context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithTimeout - Sets API server request timeout.
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.reqTimeout = t
	}
}
