package renderer

import (
	"bitbucket.org/moovie/renderer/pkg/template"
	"golang.org/x/net/context"
)

// Compiler - Components compiler.
type Compiler struct {
	Storage
}

type compilerCtxKey struct{}

// NewCompiler - Creates a new components compiler.
func NewCompiler(s Storage) *Compiler {
	return &Compiler{Storage: s}
}

// NewContext - Creates new context with compiler.
func NewContext(ctx context.Context, c *Compiler) context.Context {
	return context.WithValue(ctx, compilerCtxKey{}, c)
}

// FromContext - Retrieves compiler from context.
func FromContext(ctx context.Context) (c *Compiler, ok bool) {
	c, ok = ctx.Value(compilerCtxKey{}).(*Compiler)
	return
}

// CompileByName - Compiles a component by name.
func (compiler *Compiler) CompileByName(name string) (compiled *Compiled, err error) {
	c, err := compiler.Storage.Component(name)
	if err != nil {
		return
	}
	return compiler.Compile(c)
}

// Compile - Compiles a component.
// Expects the component to have all the required data embed or in storage.
func (compiler *Compiler) Compile(c *Component) (compiled *Compiled, err error) {
	compiled = &Compiled{Component: c}
	err = compiler.compileTo(c, compiled)
	return
}

// CompileFromStorage -
func (compiler *Compiler) CompileFromStorage(c *Component) (compiled *Compiled, err error) {
	base, err := compiler.Storage.Component(c.Name)
	if err != nil {
		return
	}
	compiled = &Compiled{Component: c}
	err = compiler.compileTo(base, compiled)
	if err != nil {
		return
	}
	err = compiler.compileTo(c, compiled)
	if err != nil {
		return
	}
	return
}

func (compiler *Compiler) compileTo(c *Component, compiled *Compiled) (err error) {
	compiled.Main, err = parseTemplate(compiler.Storage, c.Main, c.Name)
	if err != nil {
		return
	}
	compiled.Styles, err = parseTemplates(compiler.Storage, c.Styles, c.Name)
	if err != nil {
		return
	}
	compiled.Scripts, err = parseTemplates(compiler.Storage, c.Scripts, c.Name)
	if err != nil {
		return
	}
	compiled.With, err = template.ParseMap(c.With)
	if err != nil {
		return
	}
	if c.Extends != "" {
		compiled.Extends, err = compiler.CompileByName(c.Extends)
		if err != nil {
			return
		}
	}

	// Compile required components
	if len(c.Require) != 0 {
		compiled.Require = make(map[string]*Compiled)
		for name, r := range c.Require {
			compiled.Require[name], err = compiler.CompileFromStorage(r)
			if err != nil {
				return
			}
		}
	}

	return
}
