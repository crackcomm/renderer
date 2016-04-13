package middlewares

// Option - Option descriptor. Describes option and defaults.
type Option struct {
	// Name - Option name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description - Option description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Example - Example value.
	Example interface{} `json:"example,omitempty" yaml:"example,omitempty"`
	// Default - Default value.
	Default interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	// Required - Option is required to construct.
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`
}

// Options - Middleware options type.
// As read from routes configuration file.
type Options map[string]interface{}

// Clone - Clones options. Constructs map even when empty.
func (opts Options) Clone() (res Options) {
	res = make(Options)
	for key, value := range opts {
		res[key] = value
	}
	return
}
