package renderweb

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/crackcomm/renderer/middlewares"

	"gopkg.in/yaml.v2"
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
		}
		routes[route] = v
	}
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
