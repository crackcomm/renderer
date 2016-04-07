package routes

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"

	"bitbucket.org/moovie/util/template"

	"github.com/crackcomm/renderer/pkg/renderer"
	"github.com/crackcomm/renderer/pkg/web"
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
    opts:
      opts1: test_value
`

var expected = Routes{
	Route{Method: "GET", Path: "/test/:test_id"}: Handler{
		Component: renderer.Component{
			Name: "dashboard.components",
			Context: template.Context{
				"components": []interface{}{
					map[interface{}]interface{}{"name": "test1"},
					map[interface{}]interface{}{"name": "test2"},
					map[interface{}]interface{}{"name": "test3"},
				},
			},
		},
		Middlewares: []Middleware{
			{
				Name: "my_test_middleware",
				Opts: map[string]interface{}{"opts1": "test_value"},
			},
		},
	},
}

func TestRoutesUnmarshal(t *testing.T) {
	RegisterMiddleware(MiddlewareDescriptor{
		Name: "my_test_middleware",
	}, func(o Options) (web.Middleware, error) {
		return func(next xhandler.HandlerC) xhandler.HandlerC {
			return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				ctx = renderer.WithTemplateKey(ctx, "some1", "test1")
				ctx = renderer.WithTemplateKey(ctx, "some2", "test2")
				next.ServeHTTPC(ctx, w, r)
			})
		}, nil
	})

	m := make(yamlRoutes)
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

	_, err = r.Chain()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRoutesYamlMarshal(t *testing.T) {
	m := make(yamlRoutes)
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
