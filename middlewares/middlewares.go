package middlewares

import "fmt"

// Registry - Middlewares registry.
type Registry struct {
	middlewares map[string]*middleware
}

// New - Constructs new middlewares registry.
func New() *Registry {
	return &Registry{
		middlewares: make(map[string]*middleware),
	}
}

// middleware - Middleware stored in registry with additional defaults.
type middleware struct {
	Descriptor
	Constructor
	Defaults *Middleware
}

// Register - Registers a middleware constructor function.
// Panics if descriptor name is empty.
func (registry *Registry) Register(desc Descriptor, constructor Constructor) {
	if desc.Name == "" {
		panic("middleware name cannot be empty")
	}
	for _, opt := range desc.Options {
		if opt.Type == "" {
			panic(fmt.Sprintf("Middlewares %q: option %q type cannot be empty", desc.Name, opt.Name))
		}
	}
	registry.middlewares[desc.Name] = &middleware{
		Descriptor:  desc,
		Constructor: constructor,
	}
}

// Alias - Registers a middleware alias with overwritten options defaults.
func (registry *Registry) Alias(source, dest string, def *Middleware) {
	registry.middlewares[dest] = &middleware{
		Defaults:    def,
		Descriptor:  registry.middlewares[source].Descriptor,
		Constructor: registry.middlewares[source].Constructor,
	}
	registry.middlewares[dest].Descriptor.Name = dest
}

// Construct - Constructs middleware handler from name and options.
func (registry *Registry) Construct(m *Middleware) (Handler, error) {
	md, ok := registry.middlewares[m.Name]
	if !ok {
		return nil, fmt.Errorf("middleware %q doesn't exist", m.Name)
	}
	if md.Defaults != nil {
		m.SetDefaults(md.Defaults)
	}
	if err := m.Validate(md.Descriptor.Options); err != nil {
		return nil, fmt.Errorf("middleware %q: %v", md.Descriptor.Name, err)
	}
	opts, err := m.ConstructOptions()
	if err != nil {
		return nil, err
	}
	return md.Constructor(opts)
}

// Exists - Checks if middleware with given name exists in registry.
func (registry *Registry) Exists(name string) bool {
	_, ok := registry.middlewares[name]
	return ok
}

// Descriptors - Returns all registered middlewares descriptors.
func (registry *Registry) Descriptors() (desc []Descriptor) {
	for _, v := range registry.middlewares {
		desc = append(desc, v.Descriptor)
	}
	return
}

// DescriptorByName - Returns middleware descriptor by name.
func (registry *Registry) DescriptorByName(name string) (_ Descriptor, _ bool) {
	md, ok := registry.middlewares[name]
	if !ok {
		return
	}
	return md.Descriptor, true
}
