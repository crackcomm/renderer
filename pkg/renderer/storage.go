package renderer

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/golang/glog"
	"github.com/patrickmn/go-cache"
	"github.com/rjeczalik/notify"

	"bitbucket.org/moovie/renderer/pkg/template"
)

// Storage - Components storage interface.
type Storage interface {
	Component(string) (*Component, error)
	Template(string) (template.Template, error)
	Text(string) (template.Template, error)
	Close() error
}

// storage - Components storage.
type storage struct {
	dirname string
	events  chan notify.EventInfo
	cache   struct {
		components *cache.Cache
		templates  *cache.Cache
		files      *cache.Cache
	}
}

// NewStorage - Creates new components storage.
func NewStorage(dirname string, defaultExpiration time.Duration, cleanupInterval time.Duration) (Storage, error) {
	dirname, err := filepath.Abs(dirname)
	if err != nil {
		return nil, err
	}
	s := &storage{dirname: dirname}
	s.cache.components = cache.New(defaultExpiration, cleanupInterval)
	s.cache.templates = cache.New(defaultExpiration, cleanupInterval)
	s.cache.files = cache.New(defaultExpiration, cleanupInterval)
	if err := s.start(); err != nil {
		return nil, err
	}
	return s, nil
}

// Text - Returns file content as Template interface.
func (s *storage) Text(path string) (t template.Template, err error) {
	path = filepath.Join(s.dirname, path)
	body, err := s.read(path)
	if err != nil {
		return
	}
	return template.TextBytes(body), nil
}

// Template - Compiles template by file path and saves in cache.
// Returns cached template if already compiled and not changed.
func (s *storage) Template(path string) (t template.Template, err error) {
	path = filepath.Join(s.dirname, path)
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
	path := filepath.Join(s.dirname, name, "component.json")
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
	s.cache.files.Set(path, body, cache.DefaultExpiration)
	return
}

// start - Starts watching for file changes in a goroutine.
func (s *storage) start() (err error) {
	path := filepath.Join(s.dirname, "...")
	glog.Infof("[watch] start path=%q", path)

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
	for event := range s.events {
		glog.Infof("[watch] event=%q path=%q", event.Event(), event.Path())
		s.cache.files.Delete(event.Path())
		s.cache.templates.Delete(event.Path())
		s.cache.components.Delete(event.Path())
	}
}
