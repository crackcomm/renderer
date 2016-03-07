package routes

import (
	"fmt"

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

// String -
func (r Route) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.Path)
}

// Chain - Chains routes into http handler.
func (r Routes) Chain() (xhandler.HandlerC, error) {
	mux := xmux.New()
	for route, handler := range r {
		chain, err := handler.Chain()
		if err != nil {
			return nil, err
		}
		mux.HandleC(route.Method, route.Path, chain)
	}
	return mux, nil
}
