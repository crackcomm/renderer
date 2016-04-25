package middlewares

import (
	"net/http"
	"sort"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/components"
)

// DefaultRegistry - Default global registry.
var DefaultRegistry = New()

// Register - Registers a middleware constructor function.
// Panics if descriptor name is empty.
func Register(desc Descriptor, constructor Constructor) {
	DefaultRegistry.Register(desc, constructor)
}

// Alias - Registers a middleware alias with overwritten options defaults.
func Alias(source, dest string, options Options) {
	DefaultRegistry.Alias(source, dest, options)
}

// Construct - Constructs middleware handler from name and options.
func Construct(m *Middleware) (func(next xhandler.HandlerC) xhandler.HandlerC, error) {
	return DefaultRegistry.Construct(m)
}

// Exists - Checks if middleware with given name exists in registry.
func Exists(name string) bool {
	return DefaultRegistry.Exists(name)
}

// Descriptors - Returns all registered middlewares descriptors.
func Descriptors() (desc []Descriptor) {
	return DefaultRegistry.Descriptors()
}

// DescriptorByName - Returns middleware descriptor by name.
func DescriptorByName(name string) (Descriptor, bool) {
	return DefaultRegistry.DescriptorByName(name)
}

func init() {
	Register(Descriptor{
		Name:        "renderer.middlewares",
		Description: "Sets list of available middlewares descriptors in template context.",
		Context: []Option{
			{Name: "middlewares", Description: "middleware descriptors list"},
		},
	}, func(_ Options) (Handler, error) {
		descriptors := Descriptors()
		sort.Sort(byName(descriptors))
		return func(next xhandler.HandlerC) xhandler.HandlerC {
			return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				ctx = components.WithTemplateKey(ctx, "middlewares", descriptors)
				next.ServeHTTPC(ctx, w, r)
			})
		}, nil
	})
}

type byName []Descriptor

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }
