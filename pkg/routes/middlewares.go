package routes

import (
	"fmt"
	"net/http"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/pkg/renderer"
	"github.com/crackcomm/renderer/pkg/web"
)

// MiddlewareDescriptor - Web route middleware descriptor.
// Describes functionality of middleware and its options.
type MiddlewareDescriptor struct {
	// Name - Middleware name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description - Middleware description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Context - Context values set by middleware.
	Context []Option `json:"description,omitempty" yaml:"description,omitempty"`
	// Options - Options descriptors.
	Options []Option `json:"options,omitempty" yaml:"options,omitempty"`
}

// Middleware - Web route middleware.
// Name is a name of globally registered middleware.
// Its used to construct middlewares from config file.
type Middleware struct {
	// Name - Name of middleware to construct.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Options - Options used to construct middleware.
	Options Options `json:"options,omitempty" yaml:"options,omitempty"`
}

// Option - Option descriptor. Describes option and defaults.
type Option struct {
	// Name - Option name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description - Option description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Example - Example value.
	Example interface{} `json:"example,omitempty" yaml:"example,omitempty"`
	// Default - Default value.
	Default interface{} `json:"default,omitempty" yaml:"default,omitempty"`
}

// Options - Middleware options type.
// As read from routes configuration file.
type Options map[string]interface{}

// MiddlewareConstructor - Middleware constructor.
type MiddlewareConstructor func(Options) (web.Middleware, error)

type middlewaresRegistry struct {
	descriptors  map[string]MiddlewareDescriptor
	constructors map[string]MiddlewareConstructor
}

var globalMiddlewares = &middlewaresRegistry{
	descriptors:  make(map[string]MiddlewareDescriptor),
	constructors: make(map[string]MiddlewareConstructor),
}

// MiddlewaresDescriptors - Returns all registered middlewares descriptors.
func MiddlewaresDescriptors() (desc []MiddlewareDescriptor) {
	for _, v := range globalMiddlewares.descriptors {
		desc = append(desc, v)
	}
	return
}

// MiddlewareDescriptorByName - Returns middleware descriptor by name.
func MiddlewareDescriptorByName(name string) (MiddlewareDescriptor, bool) {
	desc, ok := globalMiddlewares.descriptors[name]
	return desc, ok
}

// RegisterMiddleware - Registers a middleware constructor function.
// Panics if descriptor name is empty.
func RegisterMiddleware(desc MiddlewareDescriptor, constructor MiddlewareConstructor) {
	if desc.Name == "" {
		panic("middleware name cannot be empty")
	}
	globalMiddlewares.descriptors[desc.Name] = desc
	globalMiddlewares.constructors[desc.Name] = constructor
}

// ConstructMiddleware - Constructs middleware handler from name and options.
func ConstructMiddleware(m Middleware) (web.Middleware, error) {
	md, ok := globalMiddlewares.constructors[m.Name]
	if !ok {
		return nil, fmt.Errorf("middleware %q doesn't exist", m.Name)
	}
	return md(m.Options)
}

// MiddlewareExists - Checks if middleware with given name exists in global register.
func MiddlewareExists(name string) bool {
	_, ok := globalMiddlewares.descriptors[name]
	return ok
}

// ToMiddlewareConstructor - Converts function to middleware constructor.
func ToMiddlewareConstructor(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC)) MiddlewareConstructor {
	return func(_ Options) (web.Middleware, error) {
		return web.ToMiddleware(fn), nil
	}
}

func init() {
	RegisterMiddleware(MiddlewareDescriptor{
		Name:        "web.middlewares",
		Description: "sets available middlewares descriptors in template context",
		Context: []Option{
			{Name: "middlewares", Description: "middleware descriptors list"},
		},
	}, ToMiddlewareConstructor(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		ctx = renderer.WithTemplateKey(ctx, "middlewares", MiddlewaresDescriptors())
		next.ServeHTTPC(ctx, w, r)
	}))
}
