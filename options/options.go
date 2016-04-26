package options

import "tower.pro/renderer/helpers"

// Options - Middleware options type.
// As read from routes configuration file.
type Options map[string]interface{}

// SetDefaults - Sets defaults. Deep in maps.
func (opts Options) SetDefaults(other Options) Options {
	if other == nil || len(other) == 0 {
		return opts
	}
	if opts == nil {
		opts = make(Options)
	}
	for key, value := range other {
		opts[key] = helpers.WithDefaults(opts[key], value)
	}
	return opts
}

// Clone - Clones options. Constructs map even when empty.
func (opts Options) Clone() (res Options) {
	res = make(Options)
	for key, value := range opts {
		res[key] = value
	}
	return
}
