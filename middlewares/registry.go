package middlewares

import "github.com/rs/xhandler"

// DefaultRegistry - Default global registry.
var DefaultRegistry = New()

// Register - Registers a middleware constructor function.
// Panics if descriptor name is empty.
func Register(desc Descriptor, constructor Constructor) {
	DefaultRegistry.Register(desc, constructor)
}

// Alias - Registers a middleware alias with overwritten options defaults.
func Alias(source, dest string, def *Middleware) {
	DefaultRegistry.Alias(source, dest, def)
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
