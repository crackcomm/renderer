package renderer

import (
	"bitbucket.org/moovie/renderer/pkg/template"
	"golang.org/x/net/context"
)

// Compiler - Components compiler.
type Compiler struct {
	*Storage
}

type compilerCtxKey struct{}

// NewCompiler - Creates a new components compiler.
func NewCompiler(s *Storage) *Compiler {
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

// Compile - Compiles a component.
// Expects the component to have all the required data embed or in storage.
func (compiler *Compiler) Compile(c *Component) (compiled *Compiled, err error) {
	compiled = &Compiled{Component: c}
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
	// TODO: required
	// if len(c.Require) != 0 {
	// 	compiled.Require = make(map[string]*Compiled)
	// 	for name, require := range c.Require {
	// 		compiled.Require[name], err = compiler.CompileByName(name)
	// 		if err != nil {
	// 			return
	// 		}
	// 	}
	// }
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