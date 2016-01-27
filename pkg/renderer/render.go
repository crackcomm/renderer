package renderer

import (
	"strings"

	"bitbucket.org/moovie/renderer/pkg/template"
	"github.com/PuerkitoBio/goquery"
	"github.com/flosch/pongo2"
)

// RenderHTML - Renders compiled component to html.
func RenderHTML(c *Compiled, ctx template.Context) (body string, err error) {
	res, err := Render(c, ctx)
	if err != nil {
		return
	}
	if len(res.Styles) == 0 && len(res.Scripts) == 0 {
		return res.Body, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.Body))
	if err != nil {
		return
	}

	if len(res.Styles) > 0 {
		// Find <head> element or insert if not found
		h := doc.Find("head")
		if h.Size() == 0 {
			h = doc.Find("html").PrependHtml("<head></head>")
		}

		// Insert styles into head
		for _, src := range res.Styles {
			h.AppendHtml(renderStyle(src))
		}
	}

	if len(res.Scripts) > 0 {
		// Find <body> element or insert if not found
		b := doc.Find("body")
		if b.Size() == 0 {
			b = doc.Find("html").AppendHtml("<body></body>")
		}

		// Insert scripts on the end of body tag
		for _, src := range res.Scripts {
			b.AppendHtml(renderScript(src))
		}
	}

	return doc.Html()
}

// Render - Renders compiled component.
func Render(c *Compiled, ctx template.Context) (res *Rendered, err error) {
	if ctx == nil {
		ctx = make(template.Context)
	}
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
