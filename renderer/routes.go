package renderer

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

// Route - Web route.
type Route struct {
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	Method string `json:"method,omitempty" yaml:"method,omitempty"`
}

// Routes - Routes map.
type Routes map[Route]*Handler

// String - Returns string representation of a route.
func (route Route) String() string {
	return strings.Join([]string{route.Method, route.Path}, " ")
}

// Construct - Constructs http router.
func (routes Routes) Construct(options ...Option) (xhandler.HandlerC, error) {
	// Create new router
	mux := xmux.New()

	tracing := tracingEnabled(options...)

	// Bind all routes handlers
	for route, handler := range routes {
		// Construct handler
		h, err := handler.Construct(options...)
		if err != nil {
			return nil, fmt.Errorf("%q: %v", route, err)
		}

		if tracing {
			h = routeTracing(route, h)
		}

		// Bind route handler
		mux.HandleC(route.Method, route.Path, h)
	}

	// Return handler
	return mux, nil
}

// ToStringMap - To map with string routes.
func (routes Routes) ToStringMap() (res map[string]*Handler) {
	res = make(map[string]*Handler)
	for k, v := range routes {
		res[k.String()] = v
	}
	return
}

func routeTracing(route Route, handler xhandler.HandlerC) xhandler.HandlerC {
	rs := route.String()
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		tr := trace.New(rs, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		ctx = trace.NewContext(ctx, tr)
		handler.ServeHTTPC(ctx, w, r)
		tr.Finish()
	})
}
