package middlewares

import (
	"fmt"

	"github.com/rs/xhandler"
)

// Middleware - Web route middleware.
// Name is a name of globally registered middleware.
// Its used to construct middlewares from config file.
type Middleware struct {
	// Name - Name of middleware to construct.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Options - Options used to construct middleware.
	Options Options `json:"options,omitempty" yaml:"options,omitempty"`
}

// Descriptor - Web route middleware descriptor.
// Describes functionality of middleware and its options.
type Descriptor struct {
	// Name - Middleware name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description - Middleware description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Context - Context values set by middleware.
	Context []Option `json:"description,omitempty" yaml:"description,omitempty"`
	// Options - Options descriptors.
	Options []Option `json:"options,omitempty" yaml:"options,omitempty"`
}

// Constructor - Middleware constructor.
type Constructor func(Options) (Handler, error)

// Handler - Middleware http handler.
type Handler func(next xhandler.HandlerC) xhandler.HandlerC

// SetDefaults - Sets default options and returns error if some required are missing.
func (desc Descriptor) SetDefaults(opts Options) (res Options, err error) {
	res = opts.Clone()
	for _, opt := range desc.Options {
		_, has := res[opt.Name]
		if has {
			continue
		}
		res[opt.Name] = setDefault(res[opt.Name], opt.Default)
		if opt.Required && res[opt.Name] == nil {
			return nil, fmt.Errorf("middleware %q requires %q option value", desc.Name, opt.Name)
		}
	}
	return
}

func setDefault(in, def interface{}) interface{} {
	if in == nil {
		return def
	}
	if def == nil {
		return in
	}
	return mapmerge(in, def)
}

// WithDefaults -
func (desc Descriptor) WithDefaults(opts Options) (_ Descriptor, err error) {
	desc.Options = append([]Option(nil), desc.Options...)
	for name, value := range opts {
		n, ok := getOptByName(desc.Options, name)
		if !ok {
			err = fmt.Errorf("%q has no option named %q", desc.Name, name)
			return
		}
		desc.Options[n].Default = value
	}
	return desc, nil
}

func getOptByName(opts []Option, name string) (n int, ok bool) {
	for n, opt := range opts {
		if opt.Name == name {
			return n, true
		}
	}
	return
}
