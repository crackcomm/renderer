package renderer

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"tower.pro/renderer/components"
	"tower.pro/renderer/helpers"
	"tower.pro/renderer/middlewares"
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
	for r, h := range file {
		var route Route
		route, err = parseRoute(r)
		if err != nil {
			return
		}

		if err = cleanComponent(h.Component); err != nil {
			return
		}
		for _, m := range h.Middlewares {
			if !middlewares.Exists(m.Name) {
				err = fmt.Errorf("middleware %q doesn't exist", m.Name)
				return
			}
			if err := cleanMiddleware(m); err != nil {
				return nil, fmt.Errorf("route %q: %v", r, err)
			}
		}
		routes[route] = h
	}
	return
}

func cleanComponent(c *components.Component) (err error) {
	if c == nil {
		return
	}
	c.Context, err = helpers.CleanDeepMap(c.Context)
	if err != nil {
		return
	}
	c.With, err = helpers.CleanDeepMap(c.With)
	return
}

func cleanMiddleware(m *middlewares.Middleware) (err error) {
	m.Options, err = helpers.CleanDeepMap(m.Options)
	if err != nil {
		return
	}
	m.Template, err = helpers.CleanDeepMap(m.Template)
	if err != nil {
		return
	}
	m.Context, err = helpers.CleanDeepMap(m.Context)
	return
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
