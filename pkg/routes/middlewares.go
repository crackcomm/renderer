package routes

import (
	"fmt"

	"github.com/crackcomm/renderer/pkg/web"
)

// Middleware - Web route middleware.
// Name is a name of globally registered middleware.
// Opts are options directed to this
type Middleware struct {
	Name string
	Opts Options
}

// Options - Middleware options type.
type Options map[string]interface{}

// Construct - Constructs middleware handler from name and options.
func (m Middleware) Construct() (web.Middleware, error) {
	md, ok := globalMiddlewares[m.Name]
	if !ok {
		return nil, fmt.Errorf("middleware %q doesn't exist", m.Name)
	}
	return md(m.Opts)
}

// Exists - Checks if middleware with given name exists in global register.
func (m Middleware) Exists() bool {
	_, ok := globalMiddlewares[m.Name]
	return ok
}

var globalMiddlewares = make(map[string]func(Options) (web.Middleware, error))

// RegisterMiddleware - Registers a middleware constructor function.
func RegisterMiddleware(name string, constructor func(Options) (web.Middleware, error)) {
	globalMiddlewares[name] = constructor
}
