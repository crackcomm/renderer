package web

import (
	"time"

	"github.com/rs/xhandler"
)

// NewAPI - New renderer web server API handler.
// Context should have a compiler set with `renderer.NewContext`.
//
// Default options are:
//
// 	WithTimeout(time.Second * 15),
// 	WithComponentSetter(UnmarshalFromRequest()),
//
func NewAPI(opts ...Option) xhandler.HandlerC {
	o := &options{
		reqTimeout: time.Second * 15,
	}

	for _, opt := range opts {
		opt(o)
	}

	// Create a Chain of http handlers
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
	if o.templateCtxSetter != nil {
		chain.UseC(o.templateCtxSetter)
	}
	chain.UseC(CompileFromCtx)
	chain.UseC(RenderFromCtx)

	return chain.HandlerCF(WriteRendered)
}

// Option - Sets web server handler options.
type Option func(*options)

type options struct {
	// reqTimeout - HTTP request timeout.
	reqTimeout time.Duration

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
