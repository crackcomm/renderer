package renderer

import "bitbucket.org/moovie/renderer/pkg/template"

func renderTemplates(l []template.Template, ctx template.Context) (res []string, err error) {
	for _, tmp := range l {
		r, err := template.ExecuteToString(tmp, ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return
}

func mergeComponentCtx(c *Compiled, ctx template.Context) (_ template.Context, err error) {
	// Merge with component base context
	ctx = mergeCtx(ctx, c.Context)

	// Execute component's `With` templates
	w, err := c.With.Execute(ctx)
	if err != nil {
		return
	}

	// Merge `With` into context
	ctx = mergeCtx(ctx, w)
	return ctx, nil
}

func mergeCtx(dest, source template.Context) template.Context {
	if len(source) == 0 {
		return dest
	}
	if dest == nil {
		dest = make(template.Context)
	}
	for key, value := range source {
		dest[key] = value
	}
	return dest
}

func mergeStringSlices(dest, source []string) []string {
	for _, v := range source {
		if !sliceHasString(dest, v) {
			dest = append(dest, v)
		}

	}
	return dest
}

func mergeStringMaps(dest, source map[string]string) {
	for key, value := range source {
		dest[key] = value
	}
}

func sliceHasString(slice []string, str string) bool {
	for _, v := range slice {
		if str == v {
			return true
		}
	}
	return false
}

func mergeTemplatesMap(t template.Map, m map[string]string) (_ template.Map, err error) {
	w, err := template.ParseMap(m)
	if err != nil {
		return
	}
	if t == nil {
		return w, nil
	}
	for k, v := range w {
		t[k] = v
	}
	return t, nil
}
