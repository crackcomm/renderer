package compiler

import (
	"os"
	"strings"

	"github.com/crackcomm/renderer/components"
	"github.com/crackcomm/renderer/storage"
	"github.com/crackcomm/renderer/template"
)

// Compiler - Components compiler interface.
type Compiler struct {
	*storage.Storage
}

// New - Creates a new components compiler.
func New(s *storage.Storage) *Compiler {
	return &Compiler{Storage: s}
}

// Compile - Compiles a component.
// Expects the component to have all the required data embed or in storage.
func (comp *Compiler) Compile(c *components.Component) (compiled *components.Compiled, err error) {
	compiled = &components.Compiled{Component: c}
	err = comp.compileTo(compiled, c)
	return
}

// CompileByName - Compiles a component by name.
func (comp *Compiler) CompileByName(name string) (compiled *components.Compiled, err error) {
	c, err := comp.Storage.Component(name)
	if err != nil {
		return
	}
	return comp.Compile(c)
}

// CompileFromStorage - Gets component from storage by name and merges
// with component given in argument.
func (comp *Compiler) CompileFromStorage(c *components.Component) (compiled *components.Compiled, err error) {
	// Get component from storage by name
	component, err := comp.Storage.Component(c.Name)
	if err != nil {
		return
	}

	// Compile component from storage
	compiled = &components.Compiled{Component: c}
	err = comp.compileTo(compiled, c)
	if err != nil {
		return
	}

	// Overwrite defaults with given component settings
	err = comp.compileTo(compiled, component)
	return
}

func (comp *Compiler) compileTo(compiled *components.Compiled, c *components.Component) (err error) {
	// Set defaults from base component context
	compiled.Context = compiled.Context.WithDefaults(c.Context)

	// Component base path
	base := strings.Replace(c.Name, ".", string(os.PathSeparator), -1)

	// Compile main component template if not empty
	if c.Main != "" {
		compiled.Main, err = parseTemplate(comp.Storage, c.Main, base)
		if err != nil {
			return
		}
	}

	// Parse urls and compile styles templates
	compiled.Styles, err = parseTemplates(comp.Storage, c.Styles, base, compiled.Styles)
	if err != nil {
		return
	}

	// Parse urls and compile scripts templates
	compiled.Scripts, err = parseTemplates(comp.Storage, c.Scripts, base, compiled.Scripts)
	if err != nil {
		return
	}

	// Compile `With` templates map and merge into `compiled`
	if compiled.With != nil && c.With != nil {
		compiled.With, err = compiled.With.ParseAndMerge(c.With)
	} else if c.With != nil {
		compiled.With, err = template.ParseMap(c.With)
	}
	if err != nil {
		return
	}

	// Compile a component which this one `extends`
	if c.Extends != "" {
		compiled.Extends, err = comp.CompileByName(c.Extends)
		if err != nil {
			return
		}
	}

	// Compile required components
	for name, r := range c.Require {
		req, err := comp.CompileFromStorage(&r)
		if err != nil {
			return err
		}
		if compiled.Require == nil {
			compiled.Require = make(map[string]*components.Compiled)
		}
		compiled.Require[name] = req
	}

	return
}
