package renderweb

import (
	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// New - New renderer web server API handler.
// Context should have a compiler set with `compiler.NewContext`.
// To always render HTML instead of JSON on Accept `application/json`
// use `WithAlwaysHTML()` option.
//
// Default options are:
//
// 	* WithTimeout(time.Second * 15),
// 	* WithComponentSetter(UnmarshalFromRequest()),
//
func New(opts ...Option) xhandler.HandlerC {
	return construct(constructOpts(opts...))
}

func construct(o *webOptions) xhandler.HandlerC {
	var chain xhandler.Chain
	chain.UseC(xhandler.CloseHandler)
	chain.UseC(xhandler.TimeoutHandler(o.reqTimeout))
	chain.UseC(o.componentSetter)
	chain.UseC(o.templateCtxSetter)
	for _, m := range o.middlewares {
		chain.UseC(m)
	}
	chain.UseC(CompileInContext)
	chain.UseC(RenderInContext)
	if o.alwaysHTML {
		return chain.HandlerCF(WriteRenderedHTML)
	}
	return chain.HandlerCF(WriteRendered)
}

type breakCtxKey struct{}

var breakKey = breakCtxKey{}

// Break - Sets a break in context.
func Break(ctx context.Context) context.Context {
	return context.WithValue(ctx, breakKey, true)
}

// HasBreak - Returns true if break was set in given context.
func HasBreak(ctx context.Context) bool {
	b, _ := ctx.Value(breakKey).(bool)
	return b
}
