package renderer

import (
	"fmt"
	"path/filepath"
	"strings"

	"bitbucket.org/moovie/util/template"
)

// parseTemplate - Parse template from string.
// String may be in URL format (eq. `http://...` or `file://...`).
// Or it may contain template data in format `data:template {{ here }}`.
// Or it may contain pure text data in format `text:plain data here`.
func parseTemplate(s Storage, text, baseDir string) (t template.Template, err error) {
	scheme, rest, ok := parseScheme(text)
	if !ok {
		return nil, fmt.Errorf("invalid template url %q", text)
	}

	switch scheme {
	case "template":
		return template.FromString(rest)
	case "file":
		if baseDir != "" {
			rest = filepath.Join(baseDir, rest)
		}
		return s.Template(rest)
	case "http", "https":
		return template.Text(text), nil
	case "text":
		return template.Text(rest), nil
	}

	return
}

// parseTemplates - Parses list of templates.
func parseTemplates(s Storage, texts []string, baseDir string, start []template.Template) (res []template.Template, err error) {
	res = start
	for _, text := range texts {
		t, err := parseTemplate(s, text, baseDir)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return
}

func parseScheme(text string) (scheme, rest string, ok bool) {
	i := strings.Index(text, "://")
	if i == -1 {
		return
	}
	scheme = text[:i]
	rest = text[i+3:]
	ok = true
	return
}
