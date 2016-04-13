package middlewares

import "fmt"

// Registry - Middlewares registry.
type Registry struct {
	descriptors  map[string]Descriptor
	constructors map[string]Constructor
}

// New - Constructs new middlewares registry.
func New() *Registry {
	return &Registry{
		descriptors:  make(map[string]Descriptor),
		constructors: make(map[string]Constructor),
	}
}

// Register - Registers a middleware constructor function.
// Panics if descriptor name is empty.
func (registry *Registry) Register(desc Descriptor, constructor Constructor) {
	if desc.Name == "" {
		panic("middleware name cannot be empty")
	}
	registry.descriptors[desc.Name] = desc
	registry.constructors[desc.Name] = constructor
}

// Alias - Registers a middleware alias with overwritten options defaults.
func (registry *Registry) Alias(source, dest string, options Options) {
	desc, ok := registry.descriptors[source]
	if !ok {
		panic(fmt.Sprintf("cannot register alias %q because %q was not found", dest, source))
	}

	desc, err := desc.WithDefaults(options)
	if err != nil {
		panic(err.Error())
	}
	desc.Name = dest

	registry.Register(desc, registry.constructors[source])
}

// Construct - Constructs middleware handler from name and options.
func (registry *Registry) Construct(m *Middleware) (Handler, error) {
	md, ok := registry.constructors[m.Name]
	if !ok {
		return nil, fmt.Errorf("middleware %q doesn't exist", m.Name)
	}
	desc := registry.descriptors[m.Name]
	opts, err := desc.SetDefaults(m.Options)
	if err != nil {
		return nil, err
	}
	return md(opts)
}

// Exists - Checks if middleware with given name exists in registry.
func (registry *Registry) Exists(name string) bool {
	_, ok := registry.descriptors[name]
	return ok
}

// Descriptors - Returns all registered middlewares descriptors.
func (registry *Registry) Descriptors() (desc []Descriptor) {
	for _, v := range registry.descriptors {
		desc = append(desc, v)
	}
	return
}

// DescriptorByName - Returns middleware descriptor by name.
func (registry *Registry) DescriptorByName(name string) (Descriptor, bool) {
	desc, ok := registry.descriptors[name]
	return desc, ok
}
