package renderweb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"github.com/rs/xhandler"

	"bitbucket.org/moovie/util/httputil"
	"bitbucket.org/moovie/util/stringslice"
	"bitbucket.org/moovie/util/template"

	"github.com/crackcomm/renderer/compiler"
	"github.com/crackcomm/renderer/components"
	"github.com/crackcomm/renderer/middlewares"
)

// ToMiddleware - Converts function to middleware.
func ToMiddleware(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC)) middlewares.Handler {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			fn(ctx, w, r, next)
		})
	}
}

// ComponentMiddleware - Creates a middleware that sets given component in ctx.
func ComponentMiddleware(c *components.Component) middlewares.Handler {
	return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		ctx = components.NewContext(ctx, c)
		next.ServeHTTPC(ctx, w, r)
	})
}

// UnmarshalFromRequest - Unmarshals component using `UnmarshalFromQuery` on `GET`
// method and `` on `POST` method.
var UnmarshalFromRequest = NewUnmarshalFromRequest()

// NewUnmarshalFromRequest - Unmarshals component using `UnmarshalFromQuery` on `GET`
// method and `` on `POST` method.
func NewUnmarshalFromRequest() middlewares.Handler {
	get, post := UnmarshalFromQuery("GET"), UnmarshalFromBody("POST")
	return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		if r.Method == "GET" {
			get(next).ServeHTTPC(ctx, w, r)
		} else if r.Method == "POST" {
			post(next).ServeHTTPC(ctx, w, r)
		}
	})
}

// UnmarshalFromQuery - Unmarshals component from `json` query on certain methods.
// Stores result in context to be retrieved with `components.FromContext`.
func UnmarshalFromQuery(methods ...string) middlewares.Handler {
	return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		if len(methods) != 0 && !stringslice.Contain(methods, r.Method) {
			next.ServeHTTPC(ctx, w, r)
			return
		}

		// Read component from request
		c, err := readComponent(r)
		if err != nil {
			httputil.WriteError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		// Create a context with component and move to next handler
		ctx = components.NewContext(ctx, c)
		next.ServeHTTPC(ctx, w, r)
	})
}

func readComponent(r *http.Request) (c *components.Component, err error) {
	c = new(components.Component)
	if b := r.URL.Query().Get("json"); b != "" {
		err = json.Unmarshal([]byte(b), c)
		return
	}
	c.Name = r.URL.Query().Get("name")
	if c.Name == "" {
		return nil, errors.New("no component in request")
	}
	c.Main = r.URL.Query().Get("main")
	c.Extends = r.URL.Query().Get("extends")
	if b := r.URL.Query().Get("styles"); b != "" {
		c.Styles = strings.Split(b, ",")
	}
	if b := r.URL.Query().Get("scripts"); b != "" {
		c.Scripts = strings.Split(b, ",")
	}
	if b := r.URL.Query().Get("require"); b != "" {
		c.Require = make(map[string]components.Component)
		err = json.Unmarshal([]byte(b), &c.Require)
		if err != nil {
			return
		}
	}
	if b := r.URL.Query().Get("context"); b != "" {
		c.Context = make(template.Context)
		err = json.Unmarshal([]byte(b), &c.Context)
		if err != nil {
			return
		}
	}
	if b := r.URL.Query().Get("with"); b != "" {
		c.With = make(map[string]string)
		err = json.Unmarshal([]byte(b), &c.With)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalFromBody - Unmarshals component from request bodyCompileInContext() on certain methods.
// Stores result in context to be retrieved with `components.FromContext`.
func UnmarshalFromBody(methods ...string) middlewares.Handler {
	return ToMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request, next xhandler.HandlerC) {
		if !stringslice.Contain(methods, r.Method) {
			next.ServeHTTPC(ctx, w, r)
			return
		}
		c := new(components.Component)
		err := json.NewDecoder(r.Body).Decode(c)
		if err != nil {
			httputil.WriteError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		ctx = components.NewContext(ctx, c)
		next.ServeHTTPC(ctx, w, r)
	})
}

// CompileInContext - Compiles component from context.
// Stores result in context to be retrieved with `components.FromContext`.
func CompileInContext(next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		compiled, err := compiler.Compile(ctx)
		if err != nil {
			httputil.WriteError(w, r, http.StatusExpectationFailed, fmt.Sprintf("compile error: %v", err))
			return
		}
		ctx = components.NewCompiledContext(ctx, compiled)
		next.ServeHTTPC(ctx, w, r)
	})
}

// RenderInContext - Renders compiled component from context.
// Stores result in context to be retrieved with `components.ContextRendered`.
func RenderInContext(next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		c, ok := components.CompiledFromContext(ctx)
		if !ok {
			httputil.WriteError(w, r, http.StatusBadRequest, "component not compiled")
			return
		}
		t, _ := components.TemplateContext(ctx)
		res, err := components.Render(c, t)
		if err != nil {
			httputil.WriteError(w, r, http.StatusExpectationFailed, fmt.Sprintf("render error: %v", err))
			return
		}
		ctx = components.NewRenderedContext(ctx, res)
		next.ServeHTTPC(ctx, w, r)
	})
}

// WriteRendered - Writes rendered component from context to response writer.
// Depending on `Accept` header, it will write json or plain html body.
func WriteRendered(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		WriteRenderedJSON(ctx, w, r)
	} else {
		WriteRenderedHTML(ctx, w, r)
	}
}

// WriteRenderedJSON - Writes rendered component from context to response writer.
func WriteRenderedJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	res, ok := components.RenderedFromContext(ctx)
	if !ok {
		httputil.WriteError(w, r, http.StatusBadRequest, "component not rendered")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		glog.Warningf("[api] response encode error: %v", err)
	}
}

// WriteRenderedHTML - Writes rendered component from context to response writer.
func WriteRenderedHTML(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	res, ok := components.RenderedFromContext(ctx)
	if !ok {
		httputil.WriteError(w, r, http.StatusBadRequest, "component not rendered")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(res.HTML()))
}
