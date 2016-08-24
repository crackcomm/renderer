package template

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestContext struct{}

var _ = Suite(&TestContext{})

func (s *TestContext) TestContextGet(c *C) {
	ctx := Context{"account": Context{"meta": Context{"width": 100}}}
	width := ctx.Get("account.meta.width")
	c.Assert(width, Equals, 100)
}
