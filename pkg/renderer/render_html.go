package renderer

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RenderHTML - Merges styles and scripts into HTML body.
func RenderHTML(r *Rendered) (body string, err error) {
	// Return if no styles or scripts to add.
	if len(r.Styles) == 0 && len(r.Scripts) == 0 {
		return r.Body, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r.Body))
	if err != nil {
		return
	}

	if len(r.Styles) > 0 {
		// Find <head> element or insert if not found
		h := doc.Find("head")
		if h.Size() == 0 {
			h = doc.Find("html").PrependHtml("<head></head>")
		}

		// Insert styles into head
		for _, src := range r.Styles {
			h.AppendHtml(renderStyle(src))
		}
	}

	if len(r.Scripts) > 0 {
		// Find <body> element or insert if not found
		b := doc.Find("body")
		if b.Size() == 0 {
			b = doc.Find("html").AppendHtml("<body></body>")
		}

		// Insert scripts on the end of body tag
		for _, src := range r.Scripts {
			b.AppendHtml(renderScript(src))
		}
	}

	return doc.Html()
}

func renderStyle(src string) string {
	if strings.HasPrefix(src, "://") || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, src)
	}
	return fmt.Sprintf(`<style type="text/css">%s</style>`, src)
}

func renderScript(src string) string {
	if strings.HasPrefix(src, "://") || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return fmt.Sprintf(`<script src="%s"></script>`, src)
	}
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, src)
}
