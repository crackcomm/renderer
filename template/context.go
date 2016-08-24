package template

import "strings"

// Context - Template context.
type Context map[string]interface{}

// Get - Gets value by key.
// It might deep traverse if key doesnt exist and contains a dot.
func (ctx Context) Get(key string) (_ interface{}) {
	if v, ok := ctx[key]; ok {
		return v
	}
	if !strings.Contains(key, ".") {
		return
	}
	return ctx.getDeep(strings.Split(key, ".")...)
}

func (ctx Context) getDeep(keys ...string) (_ interface{}) {
	if len(keys) == 0 {
		return ctx
	}
	first := keys[0]
	v, ok := ctx[first]
	if !ok {
		return
	}
	if len(keys) == 1 {
		return v
	}
	switch t := v.(type) {
	case Context:
		return t.getDeep(keys[1:]...)
	}
	return
}

// WithDefaults - Sets values from `source` only if previously didnt exist.
// It may return a new map if called on nil Context.
func (ctx Context) WithDefaults(source Context) Context {
	if source == nil {
		return ctx
	}
	if ctx == nil {
		ctx = make(Context)
	}
	for key, value := range source {
		if _, has := ctx[key]; !has {
			ctx[key] = value
		}
	}
	return ctx
}

// Clone - Clones context.
func (ctx Context) Clone() (res Context) {
	res = make(Context)
	for key, value := range ctx {
		res[key] = value
	}
	return
}
