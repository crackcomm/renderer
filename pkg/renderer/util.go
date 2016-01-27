package renderer

import "bitbucket.org/moovie/util/template"

// withComponentDefaults - Returns a context with component defaults set.
func withComponentDefaults(c *Compiled, ctx template.Context) (_ template.Context, err error) {
	// Set defaults from component base context
	ctx = ctx.WithDefaults(c.Context)

	// Return if no `With` to merge with
	if len(c.With) == 0 {
		return ctx, nil
	}

	// Execute component's `With` templates
	w, err := c.With.Execute(ctx)
	if err != nil {
		return
	}

	// Merge compiled `With` into context
	ctx = ctx.Merge(w)
	return ctx, nil
}
