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

func globComponents(dir string) (res map[string]*components.Component, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, ferr error) error {
		if ferr != nil || info == nil || info.IsDir() {
			return ferr
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
