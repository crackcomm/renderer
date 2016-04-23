package helpers

// WithDefaults - Returns input with defaults.
// If input is empty returns default.
// If default is empty returns input.
// If both are not empty tries to merge them.
// Merge happens when both are maps of strings and interfaces.
func WithDefaults(in, def interface{}) interface{} {
	if in == nil {
		return def
	}
	if def == nil {
		return in
	}
	return mergeDefaults(in, def)
}

// mergeDefaults - Merges two maps. Both can be one of:
// `map[string]interface{}` or `map[interface{}]interface{}`
func mergeDefaults(dest, extra interface{}) interface{} {
	switch d := dest.(type) {
	case map[string]string:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				if _, ok := d[k]; !ok {
					d[k] = v
				}
			}
			return d
		case map[string]interface{}:
			for k, v := range e {
				if _, ok := d[k]; !ok {
					d[k] = v.(string)
				}
			}
			return d
		case map[interface{}]interface{}:
			for k, v := range e {
				if _, ok := d[k.(string)]; !ok {
					d[k.(string)] = v.(string)
				}
			}
			return d
		default:
			return dest
		}
	case map[interface{}]interface{}:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				if _, ok := d[k]; !ok {
					d[k] = v
				}
			}
			return d
		case map[interface{}]interface{}:
			for k, v := range e {
				d[k] = WithDefaults(d[k], v)
			}
			return d
		case map[string]interface{}:
			for k, v := range e {
				d[k] = WithDefaults(d[k], v)
			}
			return d
		default:
			return dest
		}
	case map[string]interface{}:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				if _, ok := d[k]; !ok {
					d[k] = v
				}
			}
			return d
		case map[interface{}]interface{}:
			for k, v := range e {
				key := k.(string)
				d[key] = WithDefaults(d[key], v)
			}
			return d
		case map[string]interface{}:
			for k, v := range e {
				d[k] = WithDefaults(d[k], v)
			}
			return d
		default:
			return dest
		}
	default:
		return dest
	}
}