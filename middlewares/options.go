package middlewares

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"tower.pro/renderer/components"
	"tower.pro/renderer/helpers"
	"tower.pro/renderer/options"
	"tower.pro/renderer/template"
)

// Options - Web route middleware constructor options interface.
type Options interface {
	// Int - Gets option int value.
	Int(context.Context, *options.Option) (int, error)

	// Bool - Gets option bool value.
	Bool(context.Context, *options.Option) (bool, error)

	// String - Gets option string value.
	String(context.Context, *options.Option) (string, error)

	// Duration - Gets duration value.
	Duration(context.Context, *options.Option) (time.Duration, error)

	// Map - Gets option map value.
	Map(context.Context, *options.Option) (template.Context, error)

	// List - Gets option list value.
	List(context.Context, *options.Option) ([]interface{}, error)

	// StringList - Gets option strings list value.
	StringList(context.Context, *options.Option) ([]string, error)
}

type middlewareOpts struct {
	options  template.Context
	context  map[string]contextNode
	template template.Map
}

func constructOptions(md *Middleware) (opts *middlewareOpts, err error) {
	opts = &middlewareOpts{
		options: md.Options,
		context: make(map[string]contextNode),
	}
	opts.template, err = template.ParseMap(template.Context(md.Template))
	if err != nil {
		return
	}
	for key, value := range md.Context {
		opts.context[key] = newContextNode(value)
	}
	return
}

// Duration - Gets option int value.
func (opts *middlewareOpts) Duration(ctx context.Context, desc *options.Option) (v time.Duration, err error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return
	}
	v, ok, err := valueDuration(value)
	if err != nil {
		return
	}
	if ok {
		return
	}
	v, ok = desc.Default.(time.Duration)
	if ok {
		return
	}
	if desc.Always {
		err = fmt.Errorf("option %q value cannot be empty", desc.Name)
	}
	return
}

// Int - Gets option int value.
func (opts *middlewareOpts) Int(ctx context.Context, desc *options.Option) (v int, err error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return
	}
	v, ok, err := valueInt(value)
	if err != nil {
		return
	}
	if ok {
		return
	}
	v, ok = desc.Default.(int)
	if ok {
		return
	}
	if desc.Always {
		err = fmt.Errorf("option %q value cannot be empty", desc.Name)
	}
	return
}

// Bool - Gets option bool value.
func (opts *middlewareOpts) Bool(ctx context.Context, desc *options.Option) (v bool, err error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return
	}
	v, ok, err := valueBool(value)
	if err != nil {
		return
	}
	if ok {
		return
	}
	if v, ok = desc.Default.(bool); ok {
		return
	}
	if desc.Always {
		err = fmt.Errorf("option %q value cannot be empty", desc.Name)
	}
	return
}

// String - Gets option string value.
func (opts *middlewareOpts) String(ctx context.Context, desc *options.Option) (v string, err error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return
	}
	var ok bool
	v, ok = value.(string)
	if ok && v != "" {
		return
	}
	v, ok = desc.Default.(string)
	if ok && v != "" {
		return
	}
	if desc.Always {
		err = fmt.Errorf("option %q value cannot be empty", desc.Name)
	}
	return
}

// Map - Gets option map value.
func (opts *middlewareOpts) Map(ctx context.Context, desc *options.Option) (res template.Context, err error) {
	res = make(template.Context)
	if temp, ok := opts.template[desc.Name]; ok {
		tempctx, ok := components.TemplateContext(ctx)
		if !ok {
			return nil, errors.New("No template context set")
		}
		v, err := temp.Execute(tempctx)
		if err != nil {
			return nil, err
		}
		res = helpers.WithDefaultsMap(res, v)
	}
	if node, ok := opts.context[desc.Name]; ok {
		opt := node.Value(ctx)
		if value, ok := opt.(template.Context); ok {
			res = helpers.WithDefaultsMap(res, value)
		} else if opt != nil {
			return nil, fmt.Errorf("invalid context value for option %q: %#v", desc.Name, opt)
		}
	}
	opt := opts.options[desc.Name]
	if options, ok := opt.(template.Context); ok {
		res = helpers.WithDefaultsMap(res, options)
	} else if opt != nil {
		return nil, fmt.Errorf("invalid option value for option %q: %#v", desc.Name, opt)
	}
	res = helpers.WithDefaultsMap(res, desc.Default)
	return
}

// List - Gets option list value.
func (opts *middlewareOpts) List(ctx context.Context, desc *options.Option) (res []interface{}, err error) {
	if temp, ok := opts.template[desc.Name]; ok {
		tempctx, _ := components.TemplateContext(ctx)
		v, err := temp.Execute(tempctx)
		if err != nil {
			return nil, err
		}
		if values, ok := v.([]interface{}); ok {
			res = append(res, values...)
		}
	}
	if node, ok := opts.context[desc.Name]; ok {
		if value, ok := node.Value(ctx).([]interface{}); ok {
			res = append(res, value...)
		}
	}
	if options, ok := opts.options[desc.Name].([]interface{}); ok {
		res = append(res, options...)
	}
	if defaults, ok := desc.Default.([]interface{}); ok {
		res = append(res, defaults...)
	}
	return
}

// StringList - Gets option string list value.
func (opts *middlewareOpts) StringList(ctx context.Context, desc *options.Option) (res []string, err error) {
	list, err := opts.List(ctx, desc)
	if err != nil {
		return
	}
	for _, value := range list {
		if str, ok := value.(string); ok && str != "" {
			res = append(res, str)
		}
	}
	return
}

func (opts *middlewareOpts) getValue(ctx context.Context, desc *options.Option) (_ interface{}, _ error) {
	node, ok := opts.context[desc.Name]
	if ok {
		if v := node.Value(ctx); v != nil {
			return v, nil
		}
	}
	template, ok := opts.template[desc.Name]
	if ok {
		tempctx, _ := components.TemplateContext(ctx)
		res, err := template.Execute(tempctx)
		if err != nil {
			return nil, err
		}
		v, _ := helpers.CleanDeep(res)
		return v, nil
	}
	if v, ok := opts.options[desc.Name]; ok {
		return v, nil
	}
	if desc.DefKey != nil {
		return ctx.Value(desc.DefKey), nil
	}
	return
}

func valueInt(value interface{}) (_ int, _ bool, _ error) {
	switch t := value.(type) {
	case int:
		return t, true, nil
	case int64:
		return int(t), true, nil
	case float64:
		return int(t), true, nil
	case string:
		v, err := strconv.Atoi(t)
		if err != nil {
			return 0, false, err
		}
		return v, true, nil
	}
	return
}

func valueDuration(value interface{}) (_ time.Duration, _ bool, _ error) {
	switch t := value.(type) {
	case time.Duration:
		return t, true, nil
	case int:
		return time.Duration(t), true, nil
	case int64:
		return time.Duration(t), true, nil
	case float64:
		return time.Duration(t), true, nil
	case string:
		if t == "" {
			return
		}
		v, err := time.ParseDuration(t)
		if err != nil {
			return 0, false, err
		}
		return v, true, nil
	}
	return
}

func valueBool(value interface{}) (_ bool, _ bool, _ error) {
	switch t := value.(type) {
	case bool:
		return t, true, nil
	case string:
		if t == "" {
			return
		}
		v, err := parseBool(t)
		if err != nil {
			return false, false, err
		}
		return v, true, nil
	}
	return
}

func parseBool(str string) (v bool, err error) {
	switch str {
	case "ok", "on", "true", "yes", "1":
		return true, nil
	case "no", "off", "false", "0":
		return false, nil
	}
	return strconv.ParseBool(str)
}

type contextMap struct {
	context map[string]contextNode
}

type contextKey struct {
	key string
}

type contextValue struct {
	value interface{}
}

type contextNode interface {
	Value(context.Context) interface{}
}

func newContextMap(ctx template.Context) (res *contextMap) {
	res = &contextMap{context: make(map[string]contextNode)}
	for key, value := range ctx {
		res.context[key] = newContextNode(value)
	}
	return
}

func newContextNode(v interface{}) (res contextNode) {
	switch t := v.(type) {
	case string:
		return &contextKey{key: t}
	case template.Context:
		return newContextMap(t)
	case map[string]interface{}:
		return newContextMap(template.Context(t))
	}
	return &contextValue{value: v}
}

func (m contextMap) Value(ctx context.Context) interface{} {
	return m.Map(ctx)
}

func (m contextMap) Map(ctx context.Context) (res template.Context) {
	if len(m.context) == 0 {
		return
	}
	res = make(template.Context)
	for k, v := range m.context {
		res[k] = v.Value(ctx)
	}
	return
}

func (k *contextKey) Value(ctx context.Context) interface{} {
	return ctx.Value(k.key)
}

func (v *contextValue) Value(_ context.Context) interface{} {
	return v.value
}
