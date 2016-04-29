package renderer

import "github.com/rs/xhandler"

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
