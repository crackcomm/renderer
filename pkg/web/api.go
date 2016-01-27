package web

import (
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// NewAPI - New renderer web server API handler.
// Context should have a compiler set with `renderer.NewContext`.
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
	if o.reader == nil {
		o.reader = UnmarshalFromRequest()
	}
	chain.UseC(o.reader)
	chain.UseC(CompileFromCtx)
	chain.UseC(RenderFromCtx)

	return chain.HandlerCtx(o.ctx, xhandler.HandlerFuncC(WriteRendered))
}

// Option - Sets web server handler options.
type Option func(*options)

type options struct {
	reqTimeout time.Duration
	ctx        context.Context

	// reader - Component reader that creates context using `renderer.ComponentCtx`.
	reader Middleware

	// templateCtxSetter - Component template context setter middleware.
	// This middleware can read context from request etc.
	// It should set template context using `renderer.WithTemplateCtx`.
	templateCtxSetter Middleware
}

// WithReader - Sets component reader HTTP request middleware.
func WithReader(reader Middleware) Option {
	return func(o *options) {
		o.reader = reader
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
