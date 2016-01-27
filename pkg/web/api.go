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
	chain.UseC(UnmarshalFromQuery("GET"))
	chain.UseC(UnmarshalFromBody("POST"))
	chain.UseC(CompileFromCtx)
	chain.UseC(RenderFromCtx)

	return chain.HandlerCtx(o.ctx, xhandler.HandlerFuncC(WriteRendered))
}

// Option - Sets web server handler options.
type Option func(*options)

type options struct {
	reqTimeout time.Duration
	ctx        context.Context
}

// WithCtx - Sets API server context.
func WithCtx(ctx context.Context) Option {
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
