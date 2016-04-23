package renderweb

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"

	"github.com/crackcomm/renderer/template"

	"github.com/crackcomm/renderer/components"
	"github.com/crackcomm/renderer/middlewares"
)

var data = `
GET /test/:test_id:
  component:
    name: dashboard.components
    context:
      components:
      - name: test1
      - name: test2
      - name: test3
  middlewares:
  - name: my_test_middleware
    options:
      opts1: test_value
`

var expected = Routes{
	Route{Method: "GET", Path: "/test/:test_id"}: &Handler{
		Component: &components.Component{
			Name: "dashboard.components",
			Context: template.Context{
				"components": []interface{}{
					map[interface{}]interface{}{"name": "test1"},
					map[interface{}]interface{}{"name": "test2"},
					map[interface{}]interface{}{"name": "test3"},
				},
			},
		},
		Middlewares: []*middlewares.Middleware{
			{
				Name:    "my_test_middleware",
				Options: middlewares.Options{"opts1": "test_value"},
			},
		},
	},
}

func TestRoutesUnmarshal(t *testing.T) {
	middlewares.Register(middlewares.Descriptor{
		Name: "my_test_middleware",
	}, func(o middlewares.Options) (middlewares.Handler, error) {
		return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
			ctx = components.WithTemplateKey(ctx, "some1", "test1")
			ctx = components.WithTemplateKey(ctx, "some2", "test2")
			next.ServeHTTPC(ctx, w, r)
		}), nil
	})

	m := make(routesFile)
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		t.Fatal(err)
	}

	r, err := m.toRoutes()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(r, expected) {
		t.Errorf("Got: \n%#v\n Expected: \n%#v\n", r, expected)
	}

	t.Logf("Got: %#v", r)

	_, err = r.Construct()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRoutesYamlMarshal(t *testing.T) {
	m := make(routesFile)
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		t.Fatal(err)
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bytes.TrimSpace(d), bytes.TrimSpace([]byte(data))) {
		t.Errorf("Got: \n%s\n", d)
	}
}
