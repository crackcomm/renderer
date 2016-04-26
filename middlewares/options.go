package middlewares

import (
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

	// Map - Gets option map value.
	Map(context.Context, *options.Option) (map[string]interface{}, error)

	// Bool - Gets option bool value.
	Bool(context.Context, *options.Option) (bool, error)

	// String - Gets option string value.
	String(context.Context, *options.Option) (string, error)

	// Duration - Gets duration value.
	Duration(context.Context, *options.Option) (time.Duration, error)
}

type middlewareOpts struct {
	options  options.Options
	context  map[string]contextNode
	template template.Map
}

func constructOptions(md *Middleware) (opts *middlewareOpts, err error) {
	opts = &middlewareOpts{
		options: md.Options,
		context: make(map[string]contextNode),
	}
	for key, value := range md.Context {
		opts.context[key] = newContextNode(value)
	}
	return
}

// Map - Gets option map value.
func (opts *middlewareOpts) Map(ctx context.Context, desc *options.Option) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	if tmpl, ok := opts.template[desc.Name]; ok {
		tctx, _ := components.TemplateContext(ctx)
		v, err := tmpl.Execute(tctx)
		if err != nil {
			return nil, err
		}
		values, ok := v.(template.Context)
		if ok {
			for key, value := range values {
				res[key] = helpers.WithDefaults(res[key], value)
			}
		}
	}
	if node, ok := opts.context[desc.Name].(*contextMap); ok {
		for key, value := range node.Map(ctx) {
			res[key] = helpers.WithDefaults(res[key], value)
		}
	}
	if options, ok := opts.options[desc.Name].(map[string]interface{}); ok {
		for key, value := range options {
			res[key] = helpers.WithDefaults(res[key], value)
		}
	}
	res = helpers.WithDefaults(res, desc.Default).(map[string]interface{})
	return
}

// Duration - Gets option int value.
func (opts *middlewareOpts) Duration(ctx context.Context, desc *options.Option) (time.Duration, error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return 0, err
	}
	v, ok, err := valueDuration(value)
	if err != nil {
		return 0, err
	}
	if ok {
		return v, nil
	}
	v, _ = desc.Default.(time.Duration)
	return v, nil
}

// Int - Gets option int value.
func (opts *middlewareOpts) Int(ctx context.Context, desc *options.Option) (int, error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return 0, err
	}
	v, ok, err := valueInt(value)
	if err != nil {
		return 0, err
	}
	if ok {
		return v, nil
	}
	v, _ = desc.Default.(int)
	return v, nil
}

// Bool - Gets option bool value.
func (opts *middlewareOpts) Bool(ctx context.Context, desc *options.Option) (bool, error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return false, err
	}
	v, ok, err := valueBool(value)
	if err != nil {
		return false, err
	}
	if ok {
		return v, nil
	}
	v, _ = desc.Default.(bool)
	return v, nil
}

// String - Gets option string value.
func (opts *middlewareOpts) String(ctx context.Context, desc *options.Option) (string, error) {
	value, err := opts.getValue(ctx, desc)
	if err != nil {
		return "", err
	}
	if v, ok := value.(string); ok {
		return v, nil
	}
	v, _ := desc.Default.(string)
	return v, nil
}

func (opts *middlewareOpts) getValue(ctx context.Context, desc *options.Option) (_ interface{}, _ error) {
	node, ok := opts.context[desc.Name]
	if ok {
		return node.Value(ctx), nil
	}
	template, ok := opts.template[desc.Name]
	if ok {
		tctx, _ := components.TemplateContext(ctx)
		return template.Execute(tctx)
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
		v, err := strconv.ParseBool(t)
		if err != nil {
			return false, false, err
		}
		return v, true, nil
	}
	return
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

func newContextMap(ctx map[string]interface{}) (res *contextMap) {
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
	case map[string]interface{}:
		return newContextMap(t)
	case options.Options:
		return newContextMap(map[string]interface{}(t))
	case template.Context:
		return newContextMap(map[string]interface{}(t))
	}
	return &contextValue{value: v}
}

func (m contextMap) Value(ctx context.Context) interface{} {
	return m.Map(ctx)
}

func (m contextMap) Map(ctx context.Context) (res map[string]interface{}) {
	if len(m.context) == 0 {
		return
	}
	res = make(map[string]interface{})
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
