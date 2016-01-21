// Package compiler implements compiler of the templates and later rendering.
package compiler

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/flosch/pongo2"
	"github.com/go-fsnotify/fsnotify"
	"github.com/golang/glog"
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
		mutex: new(sync.RWMutex),
		opts:  o,
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

func (c *compiler) Render(ctx context.Context, cmp *components.Component) (res *components.Rendered, err error) {
	compiled, ok := c.getComponent(cmp.Name)
	if !ok {
		return nil, fmt.Errorf("component %s not found", cmp.Name)
	}
	return c.renderComponent(ctx, cmp, compiled)
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
	if err = c.compileAll(ctx); err != nil {
		return
	}
	if !c.opts.watch {
		return
	}
	return c.startWatch(ctx)
}

//
// 1. Render all components `require`'d by this component
// 2. Repeat for `parent` component meaning one this is `extend`-ing (if any).
// 3. Render page with required components data in context
// 4. Render `parent` if `extends` any, with rendered component as `children`
//
func (c *compiler) renderComponent(ctx context.Context, cmp *components.Component, compiled *compiledComponent) (res *components.Rendered, err error) {
	return
}

func (c *compiler) compileAll(ctx context.Context) (err error) {
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
				return fmt.Errorf("component %s compile: %v", component.Name, err)
			}
			c.mutex.Lock()
			c.components[component.Name] = compiled
			c.mutex.Unlock()
		}

	}
	return
}

func (c *compiler) compile(ctx context.Context, cmp *components.Component) (res *compiledComponent, err error) {
	glog.V(3).Infof("[compile] component: %s", cmp.Name)

	// Component main filename
	fname := resolvePath(ctx, cmp.Main)

	// Get file absolute path
	fname, err = filepath.Abs(fname)
	if err != nil {
		return
	}

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
	return
}

func (c *compiler) startWatch(ctx context.Context) (err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	for _, dir := range c.opts.dirs {
		err = watcher.Add(dir)
		if err != nil {
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-watcher.Errors:
			return err
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				err = c.compileAll(ctx)
				if err != nil {
					return
				}
			}
		}
	}
}
