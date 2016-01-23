// Package compiler implements compiler of the templates and later rendering.
package compiler

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/flosch/pongo2"
	"github.com/golang/glog"
	"github.com/rjeczalik/notify"
	"golang.org/x/net/context"

	"bitbucket.org/moovie/renderer/components"
)

// Compiler - Ccomponents compiler interface.
type Compiler interface {
	// Start - Starts the compiler.
	// Compiles all components found in defined sources.
	// It may also watch for changes (depending on options and implementation).
	Start(context.Context) error

	// Render - Renders component.
	Render(context.Context, *components.Component) (*components.Rendered, error)

	// Component - Returns component by it's name.
	// Second parameter is false if component was not found.
	Component(string) (*components.Component, bool)

	// Components - Returns list of all components compiler knows about.
	Components() []*components.Component

	// Insert - Inserts component.
	// Compiler might save it in some backend storage depending on implementation.
	// Insert(context.Context, *components.Component) error
}

// New - Creates a new compiler.
func New(opts ...Option) Compiler {
	o := new(options)
	for _, opt := range opts {
		opt(o)
	}
	return &compiler{
		components: make(map[string]*compiledComponent),
		mutex:      new(sync.RWMutex),
		opts:       o,
	}
}

// compiler - Compiler implementation.
type compiler struct {
	// opts - Compiler options.
	opts *options

	// components - Map of compiled components by their name.
	components map[string]*compiledComponent

	// mutex - Mutex for concurrent components read/write.
	mutex *sync.RWMutex

	// watchCancel - Set when watching. On call stops watching.
	watchCancel context.CancelFunc
}

// compiledComponent - Compiled component.
// Contains definition of the component along with
// compiled component template with styles and scripts.
type compiledComponent struct {
	*components.Component
	*pongo2.Template

	styles  []*pongoOrURL
	scripts []*pongoOrURL

	withTemplates map[string]*pongo2.Template
}

func (c *compiler) Component(name string) (_ *components.Component, ok bool) {
	if r, ok := c.getComponent(name); ok {
		return r.Component, true
	}
	return
}

func (c *compiler) Components() (r []*components.Component) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, cmp := range c.components {
		r = append(r, cmp.Component)
	}
	return
}

func (c *compiler) getComponent(name string) (r *compiledComponent, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	r, ok = c.components[name]
	return
}

// func (c *compiler) Insert(ctx context.Context, component *components.Component) (err error) {
// 	compiled, err := c.compile(ctx, component)
// 	if err != nil {
// 		return
// 	}
// 	c.mutex.Lock()
// 	c.components[component.Name] = compiled
// 	c.mutex.Unlock()
// 	return
// }

func (c *compiler) Start(ctx context.Context) (err error) {
	for n, dir := range c.opts.dirs {
		c.opts.dirs[n], err = filepath.Abs(dir)
		if err != nil {
			return
		}
	}
	if err = c.compileAll(ctx); err != nil {
		return
	}
	if !c.opts.watch {
		return
	}
	return c.startWatch(ctx)
}

func (c *compiler) compileAll(ctx context.Context) (err error) {
	glog.V(3).Info("[compile] all")
	for _, dir := range c.opts.dirs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		list, err := globComponents(dir)
		if err != nil {
			return err
		}

		for basedir, component := range list {
			compiled, err := c.compile(withComponentPath(ctx, basedir), component)
			if err != nil {
				return fmt.Errorf("component %q compile: %v", component.Name, err)
			}
			c.mutex.Lock()
			c.components[component.Name] = compiled
			c.mutex.Unlock()
		}

	}
	return
}

func (c *compiler) compile(ctx context.Context, cmp *components.Component) (res *compiledComponent, err error) {
	glog.V(3).Infof("[compile] name=%q", cmp.Name)

	// Component main filename
	fname := resolvePath(ctx, cmp.Main)

	// Check if filepath is allowed
	if !pathInList(fname, c.opts.dirs) {
		return nil, fmt.Errorf("template path %q is not allowed", cmp.Main)
	}

	// Compile template from file
	res = &compiledComponent{Component: cmp}
	res.Template, err = pongo2.FromFile(fname)
	if err != nil {
		return
	}

	// Compile `With` templates
	res.withTemplates, err = compileTemplatesMap(cmp.With)
	if err != nil {
		return
	}

	// Compile scripts templates
	res.scripts, err = newListPongoOrURL(ctx, c.opts.dirs, cmp.Scripts)
	if err != nil {
		return
	}

	// Compile styles templates
	res.styles, err = newListPongoOrURL(ctx, c.opts.dirs, cmp.Styles)
	if err != nil {
		return
	}

	return
}

func (c *compiler) startWatch(ctx context.Context) (err error) {
	ch := make(chan notify.EventInfo, 1)
	defer notify.Stop(ch)

	for _, dir := range c.opts.dirs {
		err = notify.Watch(fmt.Sprintf("%s/...", strings.TrimRight(dir, "/")), ch, notify.All)
		if err != nil {
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			glog.V(1).Info("[watch] event")
			err = c.compileAll(ctx)
			if err != nil {
				return
			}
		}
	}
}

func (c *compiler) Render(ctx context.Context, cmp *components.Component) (res *components.Rendered, err error) {
	compiled, ok := c.getComponent(cmp.Name)
	if !ok {
		return nil, fmt.Errorf("component %q not found", cmp.Name)
	}
	glog.V(3).Infof("[render] name=%q", cmp.Name)

	// Create temporary template context
	tempctx := make(map[string]interface{})

	// 1. Merge base context
	for key, value := range compiled.Component.Context {
		tempctx[key] = value
	}

	// 1a. Merge provided context
	for key, value := range cmp.Context {
		tempctx[key] = value
	}

	// 2. Execute templates from `compiled` base component `With`
	for key, template := range compiled.withTemplates {
		tempctx[key], err = template.Execute(pongo2.Context(tempctx))
		if err != nil {
			return
		}
	}

	// 2a. Compile and execute `With` templates from request `Component`
	templates, err := compileTemplatesMap(cmp.With)
	if err != nil {
		return
	}
	for key, template := range templates {
		tempctx[key], err = template.Execute(pongo2.Context(tempctx))
		if err != nil {
			return
		}
	}

	// Create result structure
	res = new(components.Rendered)

	// 3. Execute `compiled` (base) `required` components.
	for key, required := range compiled.Component.Require {
		r, err := c.Render(ctx, required)
		if err != nil {
			return nil, err
		}

		for link, target := range r.Links {
			if res.Links == nil {
				res.Links = make(map[string]string)
			}
			res.Links[link] = target
		}

		res.Styles = append(res.Styles, r.Styles...)
		res.Scripts = append(res.Scripts, r.Scripts...)

		tempctx[key] = r.Body
	}

	// Append compiled scripts and styles to response
	res.Scripts, err = executePongoOrURL(tempctx, compiled.scripts)
	if err != nil {
		return
	}
	res.Styles, err = executePongoOrURL(tempctx, compiled.styles)
	if err != nil {
		return
	}

	// Render scripts from request
	ptmpls, err := newListPongoOrURL(ctx, c.opts.dirs, cmp.Scripts)
	if err != nil {
		return
	}
	s, err := executePongoOrURL(tempctx, ptmpls)
	if err != nil {
		return
	}
	res.Scripts = append(res.Scripts, s...)

	// Render styles from request
	ptmpls, err = newListPongoOrURL(ctx, c.opts.dirs, cmp.Styles)
	if err != nil {
		return
	}
	s, err = executePongoOrURL(tempctx, ptmpls)
	if err != nil {
		return
	}
	res.Styles = append(res.Styles, s...)

	// 4. Render template
	res.Body, err = compiled.Template.Execute(pongo2.Context(tempctx))
	if err != nil {
		return
	}

	// 5. If `extends`: return parent with rendered template as `children` in context.
	var extends string
	if cmp.Extends == "" {
		extends = compiled.Component.Extends
	} else {
		extends = cmp.Extends
	}
	if extends == "" {
		finishRender(res)
		return
	}

	// Set rendered component as `children` in context
	// for the parent component
	tempctx["children"] = pongo2.AsSafeValue(res.Body)

	// Render parent component
	e, err := c.Render(ctx, &components.Component{
		Name:    extends,
		Context: tempctx,
	})
	if err != nil {
		return
	}

	// Add links to parent component
	if e.Links == nil {
		e.Links = res.Links
	} else {
		for link, target := range res.Links {
			e.Links[link] = target
		}
	}

	e.Styles = append(e.Styles, res.Styles...)
	e.Scripts = append(e.Scripts, res.Scripts...)
	finishRender(e)

	return e, nil
}

func finishRender(c *components.Rendered) {
	c.Styles = removeDuplicates(c.Styles)
	c.Scripts = removeDuplicates(c.Scripts)
	c.Body = replaceMultipleWhitespace(c.Body)
}

func removeDuplicates(elements []string) (result []string) {
	if len(elements) == 0 {
		return elements
	}
	done := make(map[string]bool)
	for v := range elements {
		if !done[elements[v]] {
			done[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return
}