package renderweb

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/crackcomm/renderer/helpers"
	"github.com/crackcomm/renderer/middlewares"
	"github.com/crackcomm/renderer/options"
)

// RoutesFromFile - Reads routes from yaml file.
func RoutesFromFile(filename string) (Routes, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	m := make(routesFile)
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}
	return m.toRoutes()
}

type routesFile map[string]*Handler

func (file routesFile) toRoutes() (routes Routes, err error) {
	routes = make(Routes)
	for r, v := range file {
		var route Route
		route, err = parseRoute(r)
		if err != nil {
			return
		}
		for _, m := range v.Middlewares {
			if !middlewares.Exists(m.Name) {
				err = fmt.Errorf("middleware %q doesn't exist", m.Name)
				return
			}
			if err := cleanMiddleware(m); err != nil {
				return nil, fmt.Errorf("route %q: %v", r, err)
			}
		}
		routes[route] = v
	}
	return
}

func cleanMiddleware(m *middlewares.Middleware) error {
	// Clean middleware options
	// because yaml gives us ugly maps map[interface{}]interface{}
	// which is not we want in 99.99% of cases and this makes me
	// write tons of unnecessary code to handle conversions
	// lets convert it always here, JSON doesn't have this problem
	v, ok := helpers.CleanMapDeep(map[string]interface{}(m.Options))
	if !ok {
		return fmt.Errorf("middleware %q invalid options", m.Name)
	}

	m.Options = options.Options(v.(map[string]interface{}))

	// Clean templates map
	v, ok = helpers.CleanMapDeep(map[string]interface{}(m.Template))
	if !ok {
		return fmt.Errorf("middleware %q invalid templates", m.Name)
	}
	m.Template = options.Options(v.(map[string]interface{}))

	// Clean context map
	v, ok = helpers.CleanMapDeep(map[string]interface{}(m.Context))
	if !ok {
		return fmt.Errorf("middleware %q invalid context", m.Name)
	}
	m.Context = options.Options(v.(map[string]interface{}))
	return nil
}

func parseRoute(str string) (r Route, err error) {
	s := strings.Split(str, " ")
	if len(s) != 2 {
		err = fmt.Errorf("invalid route %q", str)
		return
	}
	r.Method = s[0]
	r.Path = s[1]
	return
}
