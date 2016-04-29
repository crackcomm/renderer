package renderer

import (
	"net/http"
	"strings"

	"github.com/rs/xmux"

	"golang.org/x/net/context"
)

type requestContext struct {
	*http.Request
	context.Context
}

// NewRequestContext - Creates a context with request. Created context exposes many values.
// Context request values are available under `request` namespace (raw request).
//
//   * `request` - HTTP Request (*http.Request)
//   * `request.host` - Alias for `request.url.host` (string)
//   * `request.method` - Request Method (string)
//   * `request.remote` - Request Remote address (string)
//   * `request.header` - Request Header (http.Header)
//   * `request.header.{key}` - Request Header value (string)
//   * `request.form` - Request Form (url.Values)
//   * `request.form.{key}` - Request Form value (string)
//   * `request.url` - Request URL (*url.URL)
//   * `request.url.host` - Request URL Host (*url.URL)
//   * `request.url.query` - Request URL Query (url.Values)
//   * `request.url.query.{key}` - Request URL Query value (string)
//   * `request.url.params` - Request URL Parameters (xmux.ParamHolder)
//   * `request.url.params.{key}` - Request URL Parameter value (string)
//
func NewRequestContext(ctx context.Context, req *http.Request) context.Context {
	return &requestContext{
		Context: ctx,
		Request: req,
	}
}

func (ctx *requestContext) Value(key interface{}) interface{} {
	name, ok := key.(string)
	if !ok {
		return ctx.Context.Value(key)
	}

	if !strings.Contains(name, ".") {
		if name == "request" {
			return ctx.Request
		}
		return ctx.Context.Value(key)
	}

	split := strings.Split(name, ".")
	if split[0] != "request" {
		return ctx.Context.Value(key)
	}

	switch split[1] {
	case "host":
		return ctx.Request.Host
	case "method":
		return ctx.Request.Method
	case "remote":
		return ctx.Request.RemoteAddr
	case "url":
		return ctx.getFromURL(key, split[2:]...)
	case "header":
		return ctx.getFromHeader(split[2:]...)
	case "form":
		return ctx.getFromForm(split[2:]...)
	}

	return ctx.Context.Value(key)
}

func (ctx *requestContext) getFromURL(key interface{}, rest ...string) interface{} {
	if len(rest) == 0 {
		return ctx.Request.URL
	}
	switch rest[0] {
	case "host":
		return ctx.Request.Host
	case "path":
		return ctx.Request.URL.Path
	case "query":
		if len(rest) == 0 {
			return ctx.Request.URL.Query()
		}
		return ctx.Request.URL.Query().Get(strings.Join(rest[1:], "."))
	case "params":
		if len(rest) == 0 {
			return xmux.Params(ctx.Context)
		}
		return xmux.Params(ctx.Context).Get(strings.Join(rest[1:], "."))
	default:
		return ctx.Context.Value(key)
	}
}

func (ctx *requestContext) getFromForm(rest ...string) interface{} {
	if len(rest) == 0 {
		return ctx.Request.Form
	}
	if err := ctx.Request.ParseForm(); err != nil {
		return nil
	}
	return ctx.Request.Form.Get(strings.Join(rest, "."))
}

func (ctx *requestContext) getFromHeader(rest ...string) interface{} {
	if len(rest) == 0 {
		return ctx.Request.Header
	}
	return ctx.Request.Header.Get(strings.Join(rest, "."))
}
