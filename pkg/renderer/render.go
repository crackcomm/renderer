package renderer

import "bitbucket.org/moovie/renderer/pkg/template"

// Render - Renders compiled component.
func Render(c *Compiled, ctx template.Context) (res *Rendered, err error) {
	res = new(Rendered)
	err = renderTo(c, ctx, res, res)
	return
}

func renderTo(c *Compiled, ctx template.Context, main, res *Rendered) (err error) {
	// Merge with component context
	mergeComponentCtx(c, ctx)

	// Render required components
	for name, req := range c.Require {
		r := new(Rendered)
		err = renderTo(req, ctx, main, r)
		if err != nil {
			return
		}

		// Add resulting body to context
		ctx[name] = r.Body
	}

	// Render `Main` component template
	res.Body, err = template.ExecuteToString(c.Main, ctx)
	if err != nil {
		return
	}

	// Extend a template if any
	if c.Extends != nil {
		ctx["children"] = res.Body
		err = renderTo(c.Extends, ctx, main, res)
		if err != nil {
			return
		}
	}

	// Render component styles
	tmp, err := renderTemplates(c.Styles, ctx)
	if err != nil {
		return
	}

	// Merge component scripts into main result
	main.Styles = mergeStringSlices(main.Styles, tmp)

	// Render component scripts
	tmp, err = renderTemplates(c.Scripts, ctx)
	if err != nil {
		return
	}

	// Merge component scripts into main result
	main.Scripts = mergeStringSlices(main.Scripts, tmp)

	return
}
