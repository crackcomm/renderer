package middlewares

import (
	"net/http"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// Constructor - Middleware constructor.
type Constructor func(Options) (Handler, error)

// Handler - Middleware http handler.
type Handler func(next xhandler.HandlerC) xhandler.HandlerC

// ToHandler - Converts function to middleware.
func ToHandler(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC)) Handler {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			fn(ctx, w, r, next)
		})
	}
}
