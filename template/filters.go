package template

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/flosch/pongo2"
)

func init() {
	pongo2.RegisterFilter("json", jsonFilter)
	pongo2.RegisterFilter("json_pretty", jsonPrettyFilter)
	pongo2.RegisterFilter("humanize_time", filterHumanizeTime)
	pongo2.RegisterFilter("before_now", filterBeforeNow)
	pongo2.RegisterFilter("after_now", filterAfterNow)
	pongo2.RegisterFilter("bool", filterBool)
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

func filterBeforeNow(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, isTime := in.Interface().(time.Time)
	if !isTime {
		return nil, &pongo2.Error{
			Sender:   "filter:before_now",
			ErrorMsg: "Filter input argument must be of type 'time.Time'.",
		}
	}
	return pongo2.AsValue(t.Before(time.Now())), nil
}

func filterAfterNow(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, isTime := in.Interface().(time.Time)
	if !isTime {
		return nil, &pongo2.Error{
			Sender:   "filter:past_now",
			ErrorMsg: "Filter input argument must be of type 'time.Time'.",
		}
	}
	return pongo2.AsValue(time.Now().Before(t)), nil
}

func filterHumanizeTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, isTime := in.Interface().(time.Time)
	if !isTime {
		return nil, &pongo2.Error{
			Sender:   "filter:humanize_time",
			ErrorMsg: "Filter input argument must be of type: 'time.Time'.",
		}
	}
	return pongo2.AsValue(humanize.Time(t)), nil
}

func filterBool(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() {
		return pongo2.AsValue(false), nil
	}
	switch t := in.Interface().(type) {
	case string:
		if t == "" {
			return pongo2.AsValue(false), nil
		}
		v, err := strconv.ParseBool(t)
		if err != nil {
			return nil, &pongo2.Error{
				Sender:   "filter:bool",
				ErrorMsg: "Filter input value invalid.",
			}
		}
		return pongo2.AsValue(v), nil
	case bool:
		return pongo2.AsValue(t), nil
	}
	return nil, &pongo2.Error{
		Sender:   "filter:bool",
		ErrorMsg: "Filter input value must be of type 'bool' or 'string'.",
	}
}
