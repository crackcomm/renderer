package helpers

import "tower.pro/renderer/template"

// CleanMapDeep - Cleans map deeply. Converting all map[interface{}]interface{} to map[string]interface{}.
// OK is false when obtained invalid type in map key.
func CleanMapDeep(v interface{}) (_ interface{}, ok bool) {
	switch m := v.(type) {
	case template.Context:
		return CleanMapDeep(map[string]interface{}(m))
	case map[string]interface{}:
		for k, v := range m {
			m[k], ok = CleanMapDeep(v)
			if !ok {
				return
			}
		}
		return m, true
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range m {
			var key string
			key, ok = k.(string)
			if !ok {
				return
			}
			res[key], ok = CleanMapDeep(v)
			if !ok {
				return
			}
		}
		return res, true
	case []interface{}:
		for n, v := range m {
			m[n], ok = CleanMapDeep(v)
			if !ok {
				return
			}
		}
		return m, true
	}
	return v, true
}
