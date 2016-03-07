package renderer

// Compiler - Components compiler interface.
type Compiler interface {
	Storage

	// Compile - Compiles a component.
	// Expects the component to have all the required data embed or in storage.
	Compile(*Component) (*Compiled, error)

	// CompileByName - Compiles a component by name.
	CompileByName(string) (*Compiled, error)

	// CompileFromStorage - Gets component from storage by name and merges
	// with component given in argument.
	CompileFromStorage(*Component) (*Compiled, error)
}

// NewCompiler - Creates a new components compiler.
func NewCompiler(s Storage) Compiler {
	return &compiler{Storage: s}
}

type compiler struct {
	Storage
}

// Compile - Compiles a component.
// Expects the component to have all the required data embed or in storage.
func (comp *compiler) Compile(c *Component) (compiled *Compiled, err error) {
	compiled = &Compiled{Component: c}
	err = comp.compileTo(compiled, c)
	return
}

// CompileByName - Compiles a component by name.
func (comp *compiler) CompileByName(name string) (compiled *Compiled, err error) {
	c, err := comp.Storage.Component(name)
	if err != nil {
		return
	}
	return comp.Compile(c)
}

// CompileFromStorage - Gets component from storage by name and merges
// with component given in argument.
func (comp *compiler) CompileFromStorage(c *Component) (compiled *Compiled, err error) {
	compiled, err = comp.CompileByName(c.Name)
	if err != nil {
		return
	}
	compiled.Component = c
	err = comp.compileTo(compiled, c)
	return
}

func (comp *compiler) compileTo(compiled *Compiled, c *Component) (err error) {
	// Set defaults from base component context
	compiled.Context = compiled.Context.WithDefaults(c.Context)

	// Compile main component template if not empty
	if c.Main != "" {
		compiled.Main, err = parseTemplate(comp.Storage, c.Main, c.Name)
		if err != nil {
			return
		}
	}

	// Parse urls and compile styles templates
	compiled.Styles, err = parseTemplates(comp.Storage, c.Styles, c.Name)
	if err != nil {
		return
	}

	// Parse urls and compile scripts templates
	compiled.Scripts, err = parseTemplates(comp.Storage, c.Scripts, c.Name)
	if err != nil {
		return
	}

	// Compile `With` templates map and merge into `compiled`
	compiled.With, err = compiled.With.ParseAndMerge(c.With)
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
			compiled.Require = make(map[string]*Compiled)
		}
		compiled.Require[name] = req
	}

	return
}
