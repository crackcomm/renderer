package template

import (
	"encoding/json"

	"github.com/flosch/pongo2"
)

func init() {
	pongo2.RegisterFilter("json", jsonFilter)
	pongo2.RegisterFilter("json_pretty", jsonPrettyFilter)
}

func jsonFilter(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, e *pongo2.Error) {
	body, err := json.Marshal(in.Interface())
	if err != nil {
		return nil, &pongo2.Error{ErrorMsg: err.Error()}
	}
	return pongo2.AsValue(string(body)), nil
}

func jsonPrettyFilter(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, e *pongo2.Error) {
	body, err := json.MarshalIndent(in.Interface(), "", "  ")
	if err != nil {
		return nil, &pongo2.Error{ErrorMsg: err.Error()}
	}
	return pongo2.AsValue(string(body)), nil
}
