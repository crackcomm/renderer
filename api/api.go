// Package api implements renderer HTTP API.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"

	"golang.org/x/net/context"

	"bitbucket.org/moovie/renderer/compiler"
	"bitbucket.org/moovie/renderer/components"
)

// Handler - Compiler API handler.
// Retrieves compiler from context using `compiler.FromContext`
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c, ok := compiler.FromContext(ctx)
	if !ok {
		glog.Warning("[api] compiler not found in context")
		writeError(w, http.StatusInternalServerError, "compiler not found")
		return
	}
	if r.URL.Path != "/" {
		writeError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != "POST" {
		writeError(w, http.StatusMethodNotAllowed, fmt.Sprintf("method %q is not allowed", r.Method))
		return
	}
	cmp := new(components.Component)
	if err := json.NewDecoder(r.Body).Decode(cmp); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.Render(ctx, cmp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("compile error: %v", err))
		return
	}

	w.WriteHeader(200)

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			glog.Warningf("[api] response encode error: %v", err)
		}
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(res.Body))
	}

}

func writeError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Msg  string `json:"error_msg"`
		Code int    `json:"error_code"`
	}{
		Msg:  err,
		Code: code,
	})
}
