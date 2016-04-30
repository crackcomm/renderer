package middlewares

import (
	"errors"
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
	Map(context.Context, *options.Option) (map[string]interface{}, error)

	// List - Gets option list value.
	List(context.Context, *options.Option) ([]interface{}, error)
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

// Map - Gets option map value.
func (opts *middlewareOpts) Map(ctx context.Context, desc *options.Option) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	if temp, ok := opts.template[desc.Name]; ok {
		tempctx, ok := components.TemplateContext(ctx)
		if !ok {
			return nil, errors.New("No template context set")
		}
		v, err := temp.Execute(tempctx)
		if err != nil {
			return nil, err
		}
		if value, ok := helpers.CleanMapDeep(v); ok {
			res = setDefaults(res, value)
		}
	}
	if node, ok := opts.context[desc.Name]; ok {
		if value, ok := node.Value(ctx).(map[string]interface{}); ok {
			res = setDefaults(res, value)
		}
	}
	if options, ok := opts.options[desc.Name].(map[string]interface{}); ok {
		res = setDefaults(res, options)
	}
	res = setDefaults(res, desc.Default)
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

func setDefaults(target map[string]interface{}, defaults interface{}) map[string]interface{} {
	return helpers.WithDefaults(target, defaults).(map[string]interface{})
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
		v, _ := helpers.CleanMapDeep(res)
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
