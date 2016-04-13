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
func (r *Rendered) HTML() string {
	// Return if no styles or scripts to add.
	if len(r.Styles) == 0 && len(r.Scripts) == 0 {
		return r.Body
	}

	var extras []string
	for _, src := range r.Styles {
		extras = append(extras, renderStyle(src))
	}
	for _, src := range r.Scripts {
		extras = append(extras, renderScript(src))
	}

	return insertExtras(r.Body, extras)
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
