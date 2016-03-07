package routes

import (
	"fmt"
	"net/http"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/pkg/renderer"
	"github.com/crackcomm/renderer/pkg/web"
)

// Handler - Web route handler.
type Handler struct {
	Component   *renderer.Component `json:"component,omitempty" yaml:"component,omitempty"`
	Middlewares []Middleware        `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

// Chain - Chains handler with middlewares into working http handler.
func (h *Handler) Chain() (xhandler.HandlerC, error) {
	// Create a Chain of http handlers
	chain := new(xhandler.Chain)

	// Set component in request context
	chain.UseC(web.ComponentMiddleware(h.Component))

	// Set request middlewares
	for _, md := range h.Middlewares {
		middleware, err := md.Construct()
		if err != nil {
			return nil, err
		}
		chain.UseC(middleware)
	}

	chain.UseC(func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			ps := xmux.Params(ctx)
			for _, p := range ps {
				ctx = renderer.SetTemplateCtx(ctx, fmt.Sprintf("param_%s", p.Name), p.Value)
			}
			for k, v := range r.URL.Query() {
				if len(v) > 0 {
					ctx = renderer.SetTemplateCtx(ctx, fmt.Sprintf("querystr_%s", k), v[0])
					ctx = renderer.SetTemplateCtx(ctx, fmt.Sprintf("query_%s", k), v)
				}
			}
			r.URL.Host = r.Host
			if r.URL.Scheme == "" {
				r.URL.Scheme = "http"
			}
			ctx = renderer.SetTemplateCtx(ctx, "request", r)
			next.ServeHTTPC(ctx, w, r)
		})
	})

	// Compile component, render and write
	chain.UseC(web.CompileFromCtx)
	chain.UseC(web.RenderFromCtx)
	return chain.HandlerCF(web.WriteRenderedHTML), nil
}
