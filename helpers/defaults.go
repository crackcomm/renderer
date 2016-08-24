package helpers

import "tower.pro/renderer/template"

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

// WithDefaultsMap - With defaults but converts output to template.Context.
// Use only when you are sure you will receive template.Context result.
func WithDefaultsMap(in, def interface{}) template.Context {
	v, _ := WithDefaults(in, def).(template.Context)
	return v
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
		case template.Context:
			for k, v := range e {
				if _, ok := d[k]; !ok {
					d[k] = v.(string)
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
		default:
			return dest
		}
	case map[string]interface{}:
		return mergeDefaults(template.Context(d), extra)
	case template.Context:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				if _, ok := d[k]; !ok {
					d[k] = v
				}
			}
			return d
		case map[string]interface{}:
			var res template.Context
			if d == nil {
				res = make(template.Context)
			} else {
				res = template.Context(d)
			}
			for k, v := range e {
				d[k] = WithDefaults(d[k], v)
			}
			return res
		case template.Context:
			if d == nil {
				d = make(template.Context)
			}
			for k, v := range e {
				d[k] = WithDefaults(d[k], v)
			}
			return d
		default:
			return dest
		}
	case []interface{}:
		switch e := extra.(type) {
		case []string:
			for _, value := range e {
				d = append(d, value)
			}
			return d
		case []interface{}:
			return append(d, e...)
		default:
			return d
		}
	case []string:
		switch e := extra.(type) {
		case []interface{}:
			for _, value := range e {
				if v, ok := value.(string); ok {
					d = append(d, v)
				}
			}
			return d
		case []string:
			return append(d, e...)
		default:
			return d
		}
	default:
		return dest
	}
}
