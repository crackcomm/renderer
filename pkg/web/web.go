// Package web implements set of middlewares and API for rendering components
// in http flow.
//
// Example usage of handlers/middlewares:
//
// 	// Construct API from middlewares
// 	chain.UseC(UnmarshalFromQuery("GET"))
// 	chain.UseC(UnmarshalFromBody("POST"))
// 	chain.UseC(CompileFromCtx)
// 	chain.UseC(RenderFromCtx)
// 	http.ListenAndServe(":8080", chain.HandlerCtx(ctx, xhandler.HandlerFuncC(WriteRendered)))
//
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/moovie/renderer/pkg/renderer"
	"bitbucket.org/moovie/renderer/pkg/template"

	"github.com/golang/glog"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// UnmarshalFromQuery - Unmarshals component from `json` query on certain methods.
// Stores result in context to be retrieved with `renderer.ComponentFromCtx`.
func UnmarshalFromQuery(methods ...string) func(next xhandler.HandlerC) xhandler.HandlerC {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if !methodInList(r.Method, methods) {
				next.ServeHTTPC(ctx, w, r)
				return
			}
			body := []byte(r.URL.Query().Get("json"))
			if len(body) == 0 {
				writeError(w, http.StatusBadRequest, "no component in json query parameter")
				return
			}
			c := new(renderer.Component)
			err := json.Unmarshal(body, c)
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			ctx = renderer.ComponentCtx(ctx, c)
			next.ServeHTTPC(ctx, w, r)
		})
	}
}

// UnmarshalFromBody - Unmarshals component from request bodyCompileFromCtx() on certain methods.
// Stores result in context to be retrieved with `renderer.ComponentFromCtx`.
func UnmarshalFromBody(methods ...string) func(next xhandler.HandlerC) xhandler.HandlerC {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if !methodInList(r.Method, methods) {
				next.ServeHTTPC(ctx, w, r)
				return
			}
			c := new(renderer.Component)
			err := json.NewDecoder(r.Body).Decode(c)
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			ctx = renderer.ComponentCtx(ctx, c)
			next.ServeHTTPC(ctx, w, r)
		})
	}
}

// CompileFromCtx - Compiles component from context.
// Stores result in context to be retrieved with `renderer.ComponentFromCtx`.
func CompileFromCtx(next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		c, ok := renderer.ComponentFromCtx(ctx)
		if !ok {
			writeError(w, http.StatusBadRequest, "no component set")
			return
		}
		compiler, ok := renderer.CompilerFromCtx(ctx)
		if !ok {
			writeError(w, http.StatusInternalServerError, "compiler not found")
			return
		}
		compiled, err := compiler.CompileFromStorage(c)
		if err != nil {
			writeError(w, http.StatusExpectationFailed, fmt.Sprintf("compile error: %v", err))
			return
		}
		ctx = renderer.CompiledCtx(ctx, compiled)
		next.ServeHTTPC(ctx, w, r)
	})
}

// RenderFromCtx - Renders compiled component from context.
// Stores result in context to be retrieved with `renderer.RenderedFromCtx`.
func RenderFromCtx(next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		c, ok := renderer.CompiledFromCtx(ctx)
		if !ok {
			writeError(w, http.StatusBadRequest, "component not compiled")
			return
		}
		res, err := renderer.Render(c)
		if err != nil {
			writeError(w, http.StatusExpectationFailed, fmt.Sprintf("render error: %v", err))
			return
		}
		ctx = renderer.RenderedCtx(ctx, res)
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
	res, ok := renderer.RenderedFromCtx(ctx)
	if !ok {
		writeError(w, http.StatusBadRequest, "component not rendered")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		glog.Warningf("[api] response encode error: %v", err)
	}
}

// WriteRenderedHTML - Writes rendered component from context to response writer.
func WriteRenderedHTML(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	res, ok := renderer.RenderedFromCtx(ctx)
	if !ok {
		writeError(w, http.StatusBadRequest, "component not rendered")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	body, err := renderer.RenderHTML(res)
	if err != nil {
		http.Error(w, fmt.Sprintf("html error: %v", err), http.StatusExpectationFailed)
		return
	}
	w.Write(template.CleanWhitespaces([]byte(body)))
}

func methodInList(method string, list []string) bool {
	for _, m := range list {
		if m == method {
			return true
		}
	}
	return false
}

func writeError(w http.ResponseWriter, code int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Msg  string `json:"message,omitempty"`
		Code int    `json:"code,omitempety"`
	}{
		Msg:  err,
		Code: code,
	})
}
