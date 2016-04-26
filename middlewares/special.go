package middlewares

import (
	"net/http"
	"sort"

	"tower.pro/renderer/components"
	"tower.pro/renderer/options"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

var optMiddlewareName = &options.Option{
	ID:      "renderer.middlewares.name",
	Name:    "name",
	Type:    options.TypeString,
	Short:   "Middleware name",
	Default: "middleware",
}

var optMiddlewareDestination = &options.Option{
	ID:      "renderer.middlewares.descriptor.destination",
	Name:    "middleware",
	Type:    options.TypeDestination,
	Short:   "Middlewares descriptor",
	Default: "middleware",
}

var optMiddlewaresDestination = &options.Option{
	ID:      "renderer.middlewares.descriptor.list.destination",
	Name:    "middlewares",
	Type:    options.TypeDestination,
	Short:   "Middlewares list",
	Default: "middlewares",
}

func init() {
	Register(Descriptor{
		Name:        "renderer.middlewares.get",
		Description: "Gets middleware descriptor by name",
		Options: []*options.Option{
			optMiddlewareName,
			optMiddlewareDestination,
		},
	}, middlewareByName)

	Register(Descriptor{
		Name:        "renderer.middlewares.list",
		Description: "Sets list of available middlewares descriptors in template context.",
		Options: []*options.Option{
			optMiddlewaresDestination,
		},
	}, middlewaresList)
}

func middlewaresList(opts Options) (Handler, error) {
	descriptors := Descriptors()
	sort.Sort(byName(descriptors))

	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			reskey, _ := opts.String(ctx, optMiddlewaresDestination)
			ctx = components.WithTemplateKey(ctx, reskey, descriptors)
			next.ServeHTTPC(ctx, w, r)
		})
	}, nil
}

func middlewareByName(opts Options) (Handler, error) {
	descriptors := Descriptors()
	sort.Sort(byName(descriptors))

	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			name, err := opts.String(ctx, optMiddlewareName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			reskey, _ := opts.String(ctx, optMiddlewareDestination)

			for _, desc := range descriptors {
				if desc.Name == name {
					ctx = components.WithTemplateKey(ctx, reskey, desc)
					break
				}
			}

			next.ServeHTTPC(ctx, w, r)
		})
	}, nil
}

type byName []Descriptor

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }
