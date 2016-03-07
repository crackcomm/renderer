package renderer

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/patrickmn/go-cache"
	"github.com/rjeczalik/notify"

	"bitbucket.org/moovie/util/whitespaces"

	"bitbucket.org/moovie/util/template"
)

// Storage - Components storage interface.
type Storage interface {
	// Component - Returns component by name.
	Component(string) (*Component, error)

	// Template - Compiles template by file path and saves in cache.
	// Returns cached template if already compiled and not changed.
	Template(string) (template.Template, error)

	// Text - Returns file content as Template interface.
	Text(string) (template.Template, error)

	// Close - Destroys caches and stops watching for changes.
	Close() error
}

// NewStorage - Creates new components storage.
func NewStorage(opts ...StorageOption) (_ Storage, err error) {
	o := &storageOptions{
		cacheExpiration: 5 * time.Minute,
		cleanupInterval: 1 * time.Minute,
	}
	for _, opt := range opts {
		opt(o)
	}
	o.dirname, err = filepath.Abs(o.dirname)
	if err != nil {
		return
	}
	s := &storage{opts: o}
	s.cache.components = cache.New(o.cacheExpiration, o.cleanupInterval)
	s.cache.templates = cache.New(o.cacheExpiration, o.cleanupInterval)
	s.cache.files = cache.New(o.cacheExpiration, o.cleanupInterval)
	if o.watchingChanges {
		err = s.start()
		if err != nil {
			return
		}
	}
	return s, nil
}

// storage - Components storage.
type storage struct {
	opts *storageOptions

	events chan notify.EventInfo
	cache  struct {
		components *cache.Cache
		templates  *cache.Cache
		files      *cache.Cache
	}
}

// Text - Returns file content as Template interface.
func (s *storage) Text(path string) (t template.Template, err error) {
	path = filepath.Join(s.opts.dirname, path)
	body, err := s.read(path)
	if err != nil {
		return
	}
	return template.TextBytes(body), nil
}

// Template - Compiles template by file path and saves in cache.
// Returns cached template if already compiled and not changed.
func (s *storage) Template(path string) (t template.Template, err error) {
	path = filepath.Join(s.opts.dirname, path)
	if tmp, ok := s.cache.templates.Get(path); ok {
		return tmp.(template.Template), nil
	}
	body, err := s.read(path)
	if err != nil {
		return
	}
	t, err = template.FromBytes(body)
	if err != nil {
		return
	}
	s.cache.templates.Set(path, t, cache.DefaultExpiration)
	return
}

// Component - Returns component by name.
func (s *storage) Component(name string) (c *Component, err error) {
	path := filepath.Join(s.opts.dirname, name, "component.json")
	if tmp, ok := s.cache.components.Get(path); ok {
		return tmp.(*Component), nil
	}
	body, err := s.read(path)
	if err != nil {
		return
	}
	c = new(Component)
	err = json.Unmarshal(body, c)
	if err != nil {
		return
	}
	s.cache.components.Set(path, c, cache.DefaultExpiration)
	return
}

// Close - Destroys caches and stops watching for changes.
func (s *storage) Close() (err error) {
	if s.events != nil {
		notify.Stop(s.events)
		close(s.events)
	}
	s.cache.components.Flush()
	s.cache.templates.Flush()
	s.cache.files.Flush()
	return
}

// read - reads file content or returns cached byte array
func (s *storage) read(path string) (body []byte, err error) {
	if b, ok := s.cache.files.Get(path); ok {
		return b.([]byte), nil
	}
	body, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	if s.opts.removeWhitespace {
		body = whitespaces.Clean(body)
	}
	s.cache.files.Set(path, body, cache.DefaultExpiration)
	return
}

// start - Starts watching for file changes in a goroutine.
func (s *storage) start() (err error) {
	path := filepath.Join(s.opts.dirname, "...")
	s.events = make(chan notify.EventInfo, 1)
	err = notify.Watch(path, s.events, notify.All)
	if err != nil {
		return
	}
	go s.watch()
	return
}

// watch - Watches for changes in templates files.
// If change event is emitted, compiled template is deleted from cache.
func (s *storage) watch() {
	glog.Info("[watch] started")
	base, _ := filepath.Abs(s.opts.dirname)
	for event := range s.events {
		path := event.Path()
		if p, err := filepath.Rel(base, path); err == nil {
			path = p
		}
		path = strings.Replace(path, "\\", "/", -1)
		glog.Infof("[change] %s", path)
		s.cache.files.Delete(event.Path())
		s.cache.templates.Delete(event.Path())
		s.cache.components.Delete(event.Path())
	}
}
