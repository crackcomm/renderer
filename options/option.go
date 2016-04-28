package options

// Option - Option descriptor.
type Option struct {
	ID      string        `json:"id,omitempty"`
	Type    Type          `json:"type,omitempty"`
	Name    string        `json:"name,omitempty"`
	Short   string        `json:"short,omitempty"`
	Long    string        `json:"long,omitempty"`
	Always  bool          `json:"always,omitempty"`
	Default interface{}   `json:"default,omitempty"`
	DefKey  interface{}   `json:"def_key,omitempty"`
	OneOf   []interface{} `json:"one_of,omitempty"`
}

// Exists - Checks if option exists by name.
func Exists(opts []*Option, name string) bool {
	for _, opt := range opts {
		if opt.Name == name {
			return true
		}
	}
	return false
}

// IsEmpty - Basically works only if interface is nil or empty string.
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return t == ""
	case []interface{}:
		return len(t) == 0
	case map[string]interface{}:
		return len(t) == 0
	}
	return false
}

// Clone - Clones option descriptor.
func (desc *Option) Clone() *Option {
	return &Option{
		ID:      desc.ID,
		Type:    desc.Type,
		Name:    desc.Name,
		Short:   desc.Short,
		Long:    desc.Long,
		OneOf:   desc.OneOf,
		Always:  desc.Always,
		Default: desc.Default,
	}
}

// SetName - Sets option descriptor Name.
func (desc *Option) SetName(name string) *Option {
	desc.Name = name
	return desc
}

// SetAlways - Sets option descriptor Always parameter.
func (desc *Option) SetAlways(always bool) *Option {
	desc.Always = always
	return desc
}

// SetDefault - Sets option descriptor Default parameter.
func (desc *Option) SetDefault(v interface{}) *Option {
	desc.Default = v
	return desc
}
