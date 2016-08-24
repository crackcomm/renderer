package renderer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"

	"tower.pro/renderer/options"
	"tower.pro/renderer/template"

	"tower.pro/renderer/components"
	"tower.pro/renderer/middlewares"
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
					template.Context{"name": "test1"},
					template.Context{"name": "test2"},
					template.Context{"name": "test3"},
				},
			},
		},
		Middlewares: []*middlewares.Middleware{
			{
				Name:    "my_test_middleware",
				Options: template.Context{"opts1": "test_value"},
			},
		},
	},
}

func TestRoutesUnmarshal(t *testing.T) {
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
		gotbody, e := json.MarshalIndent(r.ToStringMap(), "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		expbody, e := json.MarshalIndent(expected.ToStringMap(), "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		t.Errorf("Got: \n%s\n Expected: \n%s\n", gotbody, expbody)
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

func init() {
	middlewares.Register(middlewares.Descriptor{
		Name: "my_test_middleware",
		Options: []*options.Option{
			{Name: "opts1", Type: options.TypeString},
		},
	}, func(o middlewares.Options) (middlewares.Handler, error) {
		return middlewares.ToHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
			ctx = components.WithTemplateKey(ctx, "some1", "test1")
			ctx = components.WithTemplateKey(ctx, "some2", "test2")
			next.ServeHTTPC(ctx, w, r)
		}), nil
	})
}
