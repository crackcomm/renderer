package renderweb

import (
	"fmt"
	"strings"

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

	// Bind all routes handlers
	for route, handler := range routes {
		// Construct handler
		h, err := handler.Construct(options...)
		if err != nil {
			return nil, fmt.Errorf("%q: %v", route, err)
		}

		// Bind route handler
		mux.HandleC(route.Method, route.Path, h)
	}

	// Return handler
	return mux, nil
}
