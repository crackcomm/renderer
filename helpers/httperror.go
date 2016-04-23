package helpers

import (
	"encoding/json"
	"net/http"
)

// Error - HTTP error. JSON serializable.
type Error struct {
	Msg  string `json:"message,omitempty"`
	Code int    `json:"code,omitempty"`
}

// WriteError - Writes plain text or JSON serialized error to response
// depending on `Accept` header. If set to `application/json` writes JSON.
// Writes plain text error otherwise.
// Returns eventual marshaler or write error.
func WriteError(w http.ResponseWriter, r *http.Request, code int, err string) error {
	e := &Error{Msg: err, Code: code}
	if r.Header.Get("Accept") == "application/json" {
		return e.WriteJSON(w)
	}
	return e.Write(w)
}

// WriteJSONError - Writes response code and JSON serialized error response.
// Sets `Content-Type` header to `application/json` and sets response code.
// Returns eventual marshaler or write error.
func WriteJSONError(w http.ResponseWriter, code int, err string) error {
	return (&Error{Msg: err, Code: code}).WriteJSON(w)
}

// WriteJSON - Writes response code and JSON serialized error to http response.
// Also sets `Content-Type` header to `application/json`.
func (e *Error) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	return json.NewEncoder(w).Encode(e)
}

// Write - Writes response code and text error to http response.
// Also sets `Content-Type` header to `application/json` and .
func (e *Error) Write(w http.ResponseWriter) (err error) {
	w.WriteHeader(e.Code)
	_, err = w.Write([]byte(e.Msg))
	return
}
