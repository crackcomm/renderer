package template

import (
	"fmt"
	"strings"
)

// Parse - Parse template from string.
// String may be in URL format (eq. `http://...` or `file://...`).
// Or it may contain template data in format `data:template {{ here }}`.
// Or it may contain pure text data in format `text:plain data here`.
func Parse(text string) (t Template, err error) {
	scheme, rest, ok := parseScheme(text)
	if !ok {
		return nil, fmt.Errorf("invalid template url %q", text)
	}

	switch scheme {
	case "template":
		return FromString(rest)
	case "file":
		return FromFile(rest)
	case "http", "https":
		return Text(text), nil
	case "text":
		return Text(rest), nil
	}

	return
}

// ParseList - Parses list of templates.
func ParseList(texts []string) (res []Template, err error) {
	for _, text := range texts {
		t, err := Parse(text)
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
