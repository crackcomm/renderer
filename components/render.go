package components

import (
	"github.com/flosch/pongo2"
	"github.com/golang/glog"

	"tower.pro/renderer/helpers"
	"tower.pro/renderer/template"
)

// Render - Renders compiled component.
// Only first template context is accepted.
// Sets source component in template context under key `source_component`.
func Render(c *Compiled, ctxs ...template.Context) (res *Rendered, err error) {
	var ctx template.Context
	if len(ctxs) == 0 || ctxs[0] == nil {
		ctx = make(template.Context)
	} else {
		ctx = ctxs[0]
	}
	ctx["source_component"] = c.Component
	res = new(Rendered)
	err = renderComponent(c, res, res, ctx)
	return
}

// renderComponent - Renders a component.
// `main` is where `Styles` and `Scripts` are inserted.
// `res` is where `Body` is inserted.
func renderComponent(c *Compiled, main, res *Rendered, ctx template.Context) (err error) {
	// Set component defaults
	ctx, err = withComponentDefaults(c, ctx)
	if err != nil {
		return
	}

	if glog.V(11) {
		glog.Infof("[render] name=%q ctx=%#v", c.Name, ctx)
	}

	// Render required components
	for name, req := range c.Require {
		r := new(Rendered)
		err = renderComponent(req, main, r, ctx)
		if err != nil {
			return
		}

		// Add resulting body to context
		ctx[name] = pongo2.AsSafeValue(r.Body)
	}

	// Render `Main` component template
	if c.Main != nil {
		res.Body, err = template.ExecuteToString(c.Main, ctx)
		if err != nil {
			return
		}
	}

	// Extend a template if any
	if c.Extends != nil {
		if res.Body != "" {
			ctx["children"] = pongo2.AsSafeValue(res.Body)
		}
		err = renderComponent(c.Extends, main, res, ctx)
		if err != nil {
			return
		}
	}

	// Render component styles and scripts
	err = renderAssets(c, main, ctx)
	if err != nil {
		return
	}

	return
}

func renderAssets(c *Compiled, res *Rendered, ctx template.Context) (err error) {
	// Render component styles
	tmp, err := template.ExecuteList(c.Styles, ctx)
	if err != nil {
		return
	}

	// Merge component scripts into result
	res.Styles = helpers.MergeUnique(res.Styles, tmp)

	// Render component scripts
	tmp, err = template.ExecuteList(c.Scripts, ctx)
	if err != nil {
		return
	}

	// Merge component scripts into result
	res.Scripts = helpers.MergeUnique(res.Scripts, tmp)
	return
}

// withComponentDefaults - Returns a context with component defaults set.
func withComponentDefaults(c *Compiled, ctx template.Context) (_ template.Context, err error) {
	// Set defaults from component base context
	ctx = ctx.WithDefaults(c.Context)

	// Return if no `With` to merge with
	if len(c.With) == 0 {
		return ctx, nil
	}

	// Execute component's `With` templates
	for key, node := range c.With {
		_, has := ctx[key]
		if has {
			continue
		}
		ctx[key], err = node.Execute(ctx)
		if err != nil {
			return
		}
	}

	return ctx, nil
}
