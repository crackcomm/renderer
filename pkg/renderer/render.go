package renderer

import (
	"encoding/json"

	"bitbucket.org/moovie/renderer/pkg/template"
	"github.com/flosch/pongo2"
	"github.com/golang/glog"
)

// Render - Renders compiled component.
// Only first template context is accepted.
func Render(c *Compiled, ctxs ...template.Context) (res *Rendered, err error) {
	res = new(Rendered)
	var ctx template.Context
	if len(ctxs) == 0 {
		ctx = make(template.Context)
	} else {
		ctx = ctxs[0]
	}
	ctx["source_component"] = c.Component
	err = renderTo(c, res, res, ctx)
	return
}

func renderTo(c *Compiled, main, res *Rendered, ctx template.Context) (err error) {
	// Merge with component context
	ctx, err = mergeComponentCtx(c, ctx)
	if err != nil {
		return
	}

	if glog.V(3) {
		b, _ := json.Marshal(ctx)
		glog.Infof("[render] name=%q ctx=%s", c.Name, b)
	}

	// Render required components
	for name, req := range c.Require {
		r := new(Rendered)
		err = renderTo(req, main, r, ctx)
		if err != nil {
			return
		}

		// Add resulting body to context
		ctx[name] = pongo2.AsSafeValue(r.Body)
	}

	// Render `Main` component template
	res.Body, err = template.ExecuteToString(c.Main, ctx)
	if err != nil {
		return
	}

	// Extend a template if any
	if c.Extends != nil {
		ctx["children"] = pongo2.AsSafeValue(res.Body)
		err = renderTo(c.Extends, main, res, ctx)
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
