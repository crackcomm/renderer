package components

import (
	"fmt"
	"strings"
)

// Rendered - Rendered component.
type Rendered struct {
	// Body - Main body of the rendered component.
	Body string `json:"body,omitempty" yaml:"body,omitempty"`

	// Styles - List of styles.
	// They can be urls or list of css styles with prefix "data:text/css;".
	Styles []string `json:"styles,omitempty" yaml:"styles,omitempty"`

	// Scripts - List of scripts.
	// They can be urls or list of js scripts with prefix "data:text/javascript;".
	Scripts []string `json:"scripts,omitempty" yaml:"scripts,omitempty"`
}

// HTML - Merges styles and scripts into HTML body.
func (r *Rendered) HTML() (html string) {
	// Return if no styles or scripts to add.
	if len(r.Styles) == 0 && len(r.Scripts) == 0 {
		return r.Body
	}
	html = insertExtras(r.Body, renderList(renderStyle, r.Styles))
	html, _ = insertBefore(html, "</html>", renderList(renderScript, r.Scripts))
	return
}

func insertExtras(html string, extras []string) (res string) {
	res, ok := insertBefore(html, "</head>", extras)
	if ok {
		return
	}
	extras = []string{"<head>", strings.Join(extras, ""), "</head>"}
	res, ok = insertBefore(html, "<body", extras)
	if ok {
		return
	}
	return strings.Join([]string{
		"<!DOCTYPE html><html lang=\"en\">",
		strings.Join(extras, ""),
		"<body>", html, "</body></html>",
	}, "")
}

func insertAfter(input, after string, extras []string) (_ string, ok bool) {
	index := strings.Index(input, after)
	if index == -1 {
		return
	}
	index = index + len(after)
	extra := strings.Join(extras, "")
	return strings.Join([]string{input[:index], extra, input[index:]}, ""), true
}

func insertBefore(input, before string, extras []string) (_ string, ok bool) {
	index := strings.Index(input, before)
	if index == -1 {
		return
	}
	extra := strings.Join(extras, "")
	return strings.Join([]string{input[:index], extra, input[index:]}, ""), true
}

func renderList(fnc func(string) string, list []string) (res []string) {
	for _, src := range list {
		res = append(res, fnc(src))
	}
	return
}

func renderStyle(src string) string {
	if hasURLPrefix(src) {
		return fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, src)
	}
	return fmt.Sprintf(`<style type="text/css">%s</style>`, src)
}

func renderScript(src string) string {
	if hasURLPrefix(src) {
		return fmt.Sprintf(`<script src="%s"></script>`, src)
	}
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, src)
}

var schemes = []string{
	"http://",
	"https://",
	"://",
}

func hasURLPrefix(src string) bool {
	if strings.HasPrefix(src, "/") && len(src) <= 255 {
		return true
	}
	return strings.HasPrefix(src, "://") || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")
}
