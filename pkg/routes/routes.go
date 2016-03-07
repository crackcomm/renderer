package routes

import (
	"fmt"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"

	"github.com/crackcomm/renderer/pkg/web"
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

// Construct - Constructs  http handler.
func (r Routes) Construct(opts ...web.Option) (xhandler.HandlerC, error) {
	mux := xmux.New()
	for route, handler := range r {
		h, err := handler.Construct(opts...)
		if err != nil {
			return nil, err
		}
		mux.HandleC(route.Method, route.Path, h)
	}
	return mux, nil
}
