package middlewares

import (
	"fmt"

	"tower.pro/renderer/options"
)

// Descriptor - Web route middleware descriptor.
// Describes functionality of middleware and its options.
type Descriptor struct {
	// Name - Middleware name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description - Middleware description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Options - Options descriptors.
	Options []*options.Option `json:"options,omitempty" yaml:"options,omitempty"`
}

// Middleware - Web route middleware.
// Name is a name of globally registered middleware.
// Its used to construct middlewares from config file.
type Middleware struct {
	// Name - Name of middleware to construct.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Options - Options used to construct middleware.
	Options options.Options `json:"options,omitempty" yaml:"options,omitempty"`
	// Context - Like options but values are context keys.
	Context options.Options `json:"context,omitempty" yaml:"context,omitempty"`
	// Template - Like options but values are templates.
	Template options.Options `json:"template,omitempty" yaml:"template,omitempty"`
}

// SetDefaults - Set options defaults.
func (middleware *Middleware) SetDefaults(other *Middleware) {
	middleware.Options = middleware.Options.SetDefaults(other.Options)
	middleware.Context = middleware.Context.SetDefaults(other.Context)
	middleware.Template = middleware.Template.SetDefaults(other.Template)
}

// Validate - Validates if options have proper types and if all required are not empty.
func (middleware *Middleware) Validate(opts []*options.Option) error {
	for _, desc := range opts {
		if desc.Default != nil {
			continue
		}
		has, valid := middleware.has(desc.Name, desc.Type)
		if has && !valid {
			return fmt.Errorf("option %q: %q value should be of type %q", desc.ID, desc.Name, desc.Type)
		}
		if !has && desc.Always {
			return fmt.Errorf("option %q: %q is required", desc.ID, desc.Name)
		}
	}
	for key := range middleware.Options {
		if !options.Exists(opts, key) {
			return fmt.Errorf("option %q was provided but doesnt exist", key)
		}
	}
	for key := range middleware.Context {
		if !options.Exists(opts, key) {
			return fmt.Errorf("option %q (context) was provided but doesnt exist", key)
		}
	}
	for key := range middleware.Template {
		if !options.Exists(opts, key) {
			return fmt.Errorf("option %q (template) was provided but doesnt exist", key)
		}
	}
	return nil
}

// ConstructOptions - Creates middleware options.
func (middleware *Middleware) ConstructOptions() (_ Options, err error) {
	return constructOptions(middleware)
}

func (middleware *Middleware) has(name string, t options.Type) (has bool, valid bool) {
	if value, ok := middleware.Options[name]; ok {
		if options.IsEmpty(value) {
			return
		}
		if t == options.TypeKey || t == options.TypeTemplate {
			return true, false
		}
		return true, options.CheckType(t, value)
	}
	if value, ok := middleware.Context[name]; ok {
		if options.IsEmpty(value) {
			return
		}
		// We cannot define type of value without context
		return true, options.IsConvertible(options.TypeKey, t)
	}
	if value, ok := middleware.Template[name]; ok {
		if options.IsEmpty(value) {
			return
		}
		return true, options.IsConvertible(options.TypeTemplate, t)
	}
	return false, false
}
