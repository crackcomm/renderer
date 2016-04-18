package middlewares

// mapmerge - Merges two maps. Both can be one of:
// `map[string]interface{}` or `map[interface{}]interface{}`
func mapmerge(dest, extra interface{}) interface{} {
	switch d := dest.(type) {
	case map[string]string:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				d[k] = v
			}
			return d
		case map[string]interface{}:
			for k, v := range e {
				d[k] = v.(string)
			}
			return d
		case map[interface{}]interface{}:
			for k, v := range e {
				d[k.(string)] = v.(string)
			}
			return d
		default:
			return dest
		}
	case map[interface{}]interface{}:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				d[k] = v
			}
			return d
		case map[interface{}]interface{}:
			for k, v := range e {
				d[k] = v
			}
			return d
		case map[string]interface{}:
			for k, v := range e {
				d[k] = v
			}
			return d
		default:
			return dest
		}
	case map[string]interface{}:
		switch e := extra.(type) {
		case map[string]string:
			for k, v := range e {
				d[k] = v
			}
			return d
		case map[interface{}]interface{}:
			for k, v := range e {
				d[k.(string)] = v
			}
			return d
		case map[string]interface{}:
			for k, v := range e {
				d[k] = v
			}
			return d
		default:
			return dest
		}
	default:
		return dest
	}
}
