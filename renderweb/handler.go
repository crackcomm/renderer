package renderweb

import (
	"net/http"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/components"
	"github.com/crackcomm/renderer/middlewares"
)

// Handler - Web route handler.
type Handler struct {
	Component   *components.Component     `json:"component,omitempty" yaml:"component,omitempty"`
	Middlewares []*middlewares.Middleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

// Construct - Constructs http handler.
func (h *Handler) Construct(opts ...Option) (xhandler.HandlerC, error) {
	// Request initialization middleware
	opts = append(opts, WithMiddleware(ToMiddleware(initMiddleware)))

	// Set component-setting middleware with handler component
	opts = append(opts, WithComponentSetter(ComponentMiddleware(h.Component)))

	// Construct handler middlewares
	for _, md := range h.Middlewares {
		middleware, err := middlewares.Construct(md)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithMiddleware(middleware))
	}

	return New(opts...), nil
}

func initMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
	ctx = NewRequestContext(ctx, r)
	ctx = components.WithTemplateKey(ctx, "request", r)
	ctx = components.WithTemplateKey(ctx, "params", xmux.Params(ctx))
	next.ServeHTTPC(ctx, w, r)
}
