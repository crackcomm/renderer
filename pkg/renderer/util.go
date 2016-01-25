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

func mergeComponentCtx(c *Compiled, ctx template.Context) (err error) {
	// Merge with component base context
	mergeCtx(ctx, c.Context)

	// Execute component's `With` templates
	w, err := c.With.Execute(ctx)
	if err != nil {
		return
	}

	// Merge `With` into context
	mergeCtx(ctx, w)
	return
}

func mergeCtx(dest, source template.Context) {
	for key, value := range source {
		dest[key] = value
	}
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
