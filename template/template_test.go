package template

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestParseScheme(c *C) {
	tests := []struct {
		from, scheme, rest string
		fail               bool
	}{
		{from: "template://component.html", scheme: "template", rest: "component.html"},
		{from: "text://component.html", scheme: "text", rest: "component.html"},
		{from: "http://component.html", scheme: "http", rest: "component.html"},
		{from: "text://", scheme: "text"},
		{from: "fail test", fail: true},
	}

	for _, test := range tests {
		scheme, rest, ok := parseScheme(test.from)
		c.Logf("%s: %q, %q, %v", test.from, scheme, rest, ok)
		c.Assert(scheme, Equals, test.scheme)
		c.Assert(rest, Equals, test.rest)
		c.Assert(ok, Equals, !test.fail)
	}
}
