package renderer

import (
	"net/http"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"tower.pro/renderer/components"
	"tower.pro/renderer/middlewares"
)

// Handler - Web route handler.
type Handler struct {
	Component   *components.Component     `json:"component,omitempty" yaml:"component,omitempty"`
	Middlewares []*middlewares.Middleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

// Construct - Constructs http handler.
func (h *Handler) Construct(opts ...Option) (xhandler.HandlerC, error) {
	// Request initialization middleware
	opts = append(opts, WithMiddleware(initMiddleware))

	// Set component-setting middleware with handler component
	opts = append(opts, WithComponentSetter(ComponentMiddleware(h.Component)))

	// Check if tracing is enabled
	tracing := tracingEnabled(opts...)

	// Construct handler middlewares
	for _, md := range h.Middlewares {
		middleware, err := middlewares.Construct(md)
		if err != nil {
			return nil, err
		}
		if tracing {
			middleware = tracingMiddleware(md, middleware)
		}
		opts = append(opts, WithMiddleware(middleware))
	}

	return New(opts...), nil
}

var initMiddleware = middlewares.ToHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
	ctx = NewRequestContext(ctx, r)
	ctx = components.WithTemplateKey(ctx, "request", r)
	ctx = components.WithTemplateKey(ctx, "params", xmux.Params(ctx))
	next.ServeHTTPC(ctx, w, r)
})

// tracingMiddleware - Tracing for middlewares.
func tracingMiddleware(md *middlewares.Middleware, handler middlewares.Handler) middlewares.Handler {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		h := handler(next)
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			tr, _ := trace.FromContext(ctx)
			tr.LazyPrintf("%s", md.Name)
			h.ServeHTTPC(ctx, w, r)
		})
	}
}

func tracingEnabled(opts ...Option) bool {
	o := new(webOptions)
	for _, opt := range opts {
		opt(o)
	}
	return o.tracing
}
