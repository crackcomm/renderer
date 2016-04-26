package compiler

import (
	"fmt"
	"path/filepath"
	"strings"

	"tower.pro/renderer/storage"
	"tower.pro/renderer/template"
)

// parseTemplate - Parse template from string.
// String may be in URL format (eq. `http://...` or `file://...`).
// Or it may contain template data in format `data:template {{ here }}`.
// Or it may contain pure text data in format `text:plain data here`.
func parseTemplate(s *storage.Storage, text, baseDir string) (t template.Template, err error) {
	scheme, rest, ok := parseScheme(text)
	if !ok {
		return nil, fmt.Errorf("missing scheme from template URL: %q", text)
	}

	switch scheme {
	case "template":
		return template.FromString(rest)
	case "file":
		if baseDir != "" {
			rest = filepath.Join(baseDir, rest)
		}
		t, err = s.Template(rest)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", text, err)
		}
		return
	case "file+text":
		if baseDir != "" {
			rest = filepath.Join(baseDir, rest)
		}
		t, err = s.Text(rest)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", text, err)
		}
		return
	case "http", "https":
		return template.Text(text), nil
	case "text":
		return template.Text(rest), nil
	}

	return
}

// parseTemplates - Parses list of templates.
func parseTemplates(s *storage.Storage, texts []string, baseDir string, start []template.Template) (res []template.Template, err error) {
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
	if strings.HasPrefix(text, "/") {
		return "http", text, true
	}
	i := strings.Index(text, "://")
	if i == -1 {
		return
	}
	scheme = text[:i]
	rest = text[i+3:]
	ok = true
	return
}
