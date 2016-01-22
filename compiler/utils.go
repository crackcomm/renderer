package compiler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/moovie/renderer/components"
	"github.com/flosch/pongo2"
	"golang.org/x/net/context"
)

type pongoOrURL struct {
	template *pongo2.Template
	url      string
}

func executePongoOrURL(ctx pongo2.Context, input []*pongoOrURL) (res []string, err error) {
	for _, t := range input {
		body, err := t.execute(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, body)
	}
	return
}

func newListPongoOrURL(ctx context.Context, allowedDirs []string, list []string) (res []*pongoOrURL, err error) {
	for _, fname := range list {
		r, err := newPongoOrURL(ctx, allowedDirs, fname)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return
}

func pongoFromData(ctx context.Context, v string) (*pongoOrURL, error) {
	v = strings.TrimSpace(v[strings.Index(v, ";")+1:])
	template, err := pongo2.FromString(v)
	if err != nil {
		return nil, err
	}
	return &pongoOrURL{template: template}, nil
}

func newPongoOrURL(ctx context.Context, allowedDirs []string, v string) (*pongoOrURL, error) {
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return &pongoOrURL{url: v}, nil
	}
	if strings.HasPrefix(v, "data:") && strings.Index(v, ";") > 0 {
		return pongoFromData(ctx, v)
	}
	path := resolvePath(ctx, v)
	if !pathInList(path, allowedDirs) {
		return nil, fmt.Errorf("template path %q is not allowed", v)
	}
	template, err := pongo2.FromFile(path)
	if err != nil {
		return nil, err
	}
	return &pongoOrURL{template: template}, nil
}

func (t *pongoOrURL) execute(ctx pongo2.Context) (res string, err error) {
	if t.template == nil {
		return t.url, nil
	}
	res, err = t.template.Execute(ctx)
	if err != nil {
		return
	}
	res = replaceMultipleWhitespace(res)
	return
}

func globComponents(dir string) (res map[string]*components.Component, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, ferr error) error {
		if ferr != nil || info == nil || info.IsDir() {
			return ferr
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		dirname, fname := filepath.Split(path)
		if fname != "component.json" {
			return nil
		}

		if res == nil {
			res = make(map[string]*components.Component)
		}

		body, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		cmp := new(components.Component)
		if err := json.Unmarshal(body, cmp); err != nil {
			return err
		}

		res[dirname] = cmp
		return nil
	})
	return
}

// withComponentPath - Creates a new context with component path.
func withComponentPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, "component.path", path)
}

// resolvePath - Resolves component path using component path from context.
func resolvePath(ctx context.Context, path string) string {
	if base, ok := ctx.Value("component.path").(string); ok {
		return filepath.Join(base, path)
	}
	return path
}

// compileTemplatesMap - Compiles a map of templates
func compileTemplatesMap(input map[string]string) (result map[string]*pongo2.Template, err error) {
	if len(input) == 0 {
		return
	}
	result = make(map[string]*pongo2.Template)
	for key, value := range input {
		result[key], err = pongo2.FromString(value)
		if err != nil {
			return nil, fmt.Errorf("%q: %v", key, err)
		}
	}
	return
}

// pathInList - Checks if given path is contained in one of directories on the list.
func pathInList(path string, list []string) bool {
	for _, dir := range list {
		if path != dir && strings.HasPrefix(path, dir) {
			return true
		}
	}
	return false
}

//
// below some code
// imported from "github.com/tdewolff/parse"
//

var whitespaceTable = [256]bool{
	// ASCII
	false, false, false, false, false, false, false, false,
	false, true, true, false, true, true, false, false, // tab, new line, form feed, carriage return
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	true, false, false, false, false, false, false, false, // space
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	// non-ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}

// isWhitespace returns true for space, \n, \r, \t, \f.
func isWhitespace(c byte) bool {
	return whitespaceTable[c]
}

// replaceMultipleWhitespace replaces character series of space, \n, \t, \f, \r into a single space or newline (when the serie contained a \n or \r).
func replaceMultipleWhitespace(input string) string {
	b := []byte(input)
	j := 0
	prevWS := false
	hasNewline := false
	for i := 0; i < len(b); i++ {
		c := b[i]
		if isWhitespace(c) {
			prevWS = true
			if c == '\n' || c == '\r' {
				hasNewline = true
			}
		} else {
			if prevWS {
				prevWS = false
				if hasNewline {
					hasNewline = false
				} else {
					b[j] = ' '
					j++
				}
			}
			b[j] = b[i]
			j++
		}
	}
	if prevWS {
		if hasNewline {
			b[j] = '\n'
		} else {
			b[j] = ' '
		}
		j++
	}
	return string(b[:j])
}
