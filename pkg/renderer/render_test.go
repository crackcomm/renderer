package renderer

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bitbucket.org/moovie/renderer/pkg/template"

	. "gopkg.in/check.v1"
)

// TODO; test `extends`

type renderTest struct {
	c   []byte // JSON encoded component to render
	b   []byte // JSON encoded result expected
	ctx template.Context
}

var renderTests = []renderTest{
	{
		ctx: template.Context{
			"title":   "Test #1",
			"color":   "#fff",
			"message": "test message",
		},
		c: []byte(`{
		  "name": "example.root",
		  "main": "template://<h1>{{title}}</h1>",
		  "styles": [
		    "template://h1 { color: {{ color }}; }",
		    "text://some text here"
		  ],
		  "scripts": [
		    "template://console.log('{{ message }}');",
		    "text://console.log('{{ it_shouldnt_compile }}');"
		  ]
		}`),
		b: []byte(`{
		  "body": "\u003ch1\u003eTest #1\u003c/h1\u003e",
		  "styles": [
		    "h1 { color: #fff; }",
		    "some text here"
		  ],
		  "scripts": [
		    "console.log('test message');",
		    "console.log('{{ it_shouldnt_compile }}');"
		  ]
		}`),
	},
}

func Test(t *testing.T) { TestingT(t) }

type RenderSuite struct{}

var _ = Suite(&RenderSuite{})

func (suite *RenderSuite) TestRender(c *C) {
	s, err := NewStorage(filepath.Join(os.TempDir()), time.Minute, time.Minute)
	c.Check(err, IsNil)
	c.Check(s, NotNil)

	compiler := NewCompiler(s)
	c.Check(compiler, NotNil)

	for _, rt := range renderTests {
		component := new(Component)
		c.Check(json.Unmarshal(rt.c, component), IsNil)

		compiled, err := compiler.Compile(component)
		c.Check(err, IsNil)
		c.Check(compiled, NotNil)

		res, err := Render(compiled, rt.ctx)
		c.Check(res, NotNil)
		c.Check(err, IsNil)

		// LOG
		b, _ := json.MarshalIndent(res, "", "  ")
		log.Printf("%s", b)
		log.Printf("%s", res.Body)
		// END LOG

		e := new(Rendered)
		c.Check(json.Unmarshal(rt.b, e), IsNil)
		c.Check(res.Body, Equals, e.Body)

	}
}
