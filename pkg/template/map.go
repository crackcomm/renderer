package template

// Map - Templates map.
type Map map[string]Template

// ParseMap - Parses a map of templates.
func ParseMap(input map[string]string) (res Map, err error) {
	if len(input) == 0 {
		return
	}
	res = make(Map)
	for name, template := range input {
		res[name], err = FromString(template)
		if err != nil {
			return
		}
	}
	return
}

// Execute - Execute templates map.
func (m Map) Execute(input Context) (output Context, err error) {
	output = make(Context)
	for name, template := range m {
		output[name], err = ExecuteToString(template, input)
	}
	return
}
