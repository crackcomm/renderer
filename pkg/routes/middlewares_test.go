package routes

import (
	"net/http"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/pkg/renderer"
	"github.com/crackcomm/renderer/pkg/web"
)

func ExampleRegisterMiddleware() {
	RegisterMiddleware("context_request", func(_ Options) (web.Middleware, error) {
		return func(next xhandler.HandlerC) xhandler.HandlerC {
			return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				ctx = renderer.SetTemplateCtx(ctx, "template_key", "example value")
				next.ServeHTTPC(ctx, w, r)
			})
		}, nil
	})
}
