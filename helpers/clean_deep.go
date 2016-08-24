package helpers

import (
	"fmt"

	"tower.pro/renderer/template"
)

// CleanDeep - Cleans map deeply. Converting all map[interface{}]interface{} to map[string]interface{}.
// OK is false when obtained invalid type in map key.
func CleanDeep(v interface{}) (_ interface{}, err error) {
	var ok bool
	switch m := v.(type) {
	case template.Context:
		for k, v := range m {
			m[k], err = CleanDeep(v)
			if err != nil {
				err = fmt.Errorf("key %q: %v", k, err)
				return
			}
		}
		return m, nil
	case map[string]interface{}:
		return CleanDeep(template.Context(m))
	case map[interface{}]interface{}:
		res := make(template.Context)
		for k, v := range m {
			var key string
			key, ok = k.(string)
			if !ok {
				err = fmt.Errorf("key %#v is of invalid type %T", k, k)
				return
			}
			res[key], err = CleanDeep(v)
			if err != nil {
				err = fmt.Errorf("key %q: %v", key, err)
				return
			}
		}
		return res, nil
	case []interface{}:
		for n, v := range m {
			m[n], err = CleanDeep(v)
			if err != nil {
				err = fmt.Errorf("slice key %d: %v", n, err)
				return
			}
		}
		return m, nil
	}
	return v, nil
}

// CleanDeepMap - Use only when you are sure output is template.Context.
func CleanDeepMap(v interface{}) (_ template.Context, err error) {
	res, err := CleanDeep(v)
	if err != nil {
		return
	}
	ctx, _ := res.(template.Context)
	return ctx, nil
}
