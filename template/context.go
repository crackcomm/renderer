package template

// Context - Template context.
type Context map[string]interface{}

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

// Merge - Merges `source` into current context.
// It may return a new map if called on nil Context.
func (ctx Context) Merge(source Context) Context {
	if source == nil {
		return ctx
	}
	if ctx == nil {
		ctx = make(Context)
	}
	for key, value := range source {
		ctx[key] = value
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
