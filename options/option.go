package options

// Option - Option descriptor.
type Option struct {
	ID      string
	Type    Type
	Name    string
	Short   string
	Long    string
	OneOf   []interface{}
	Always  bool
	Default interface{}
	DefKey  interface{}
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
