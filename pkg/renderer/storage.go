package renderer

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rjeczalik/notify"

	"bitbucket.org/moovie/renderer/pkg/template"
)

// Storage - Components storage.
type Storage struct {
	dirname string
	events  chan notify.EventInfo
	cache   struct {
		components *cache.Cache
		templates  *cache.Cache
		files      *cache.Cache
	}
}

// NewStorage - Creates new components storage.
func NewStorage(dirname string, defaultExpiration time.Duration, cleanupInterval time.Duration) (s *Storage, err error) {
	s = &Storage{dirname: dirname}
	s.cache.components = cache.New(defaultExpiration, cleanupInterval)
	s.cache.templates = cache.New(defaultExpiration, cleanupInterval)
	s.cache.files = cache.New(defaultExpiration, cleanupInterval)
	err = s.start()
	return
}

// Text - Returns file content as Template interface.
func (s *Storage) Text(name string) (t template.Template, err error) {
	body, err := s.read(name)
	if err != nil {
		return
	}
	return template.TextBytes(body), nil
}

// Template - Compiles template by filename and saves in cache.
// Returns cached template if already compiled and not changed.
func (s *Storage) Template(name string) (t template.Template, err error) {
	if tmp, ok := s.cache.templates.Get(name); ok {
		return tmp.(template.Template), nil
	}
	body, err := s.read(name)
	if err != nil {
		return
	}
	t, err = template.FromBytes(body)
	if err != nil {
		return
	}
	s.cache.templates.Set(name, t, cache.DefaultExpiration)
	return
}

// Component - Returns component by name.
func (s *Storage) Component(name string) (c *Component, err error) {
	if tmp, ok := s.cache.components.Get(name); ok {
		return tmp.(*Component), nil
	}
	dirname := filepath.Join(s.dirname, name)
	body, err := s.read(filepath.Join(dirname, "component.json"))
	if err != nil {
		return
	}
	c = new(Component)
	err = json.Unmarshal(body, c)
	if err != nil {
		return
	}
	s.cache.components.Set(name, c, cache.DefaultExpiration)
	return
}

// Close - Destroys caches and stops watching for changes.
func (s *Storage) Close() (err error) {
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
func (s *Storage) read(name string) (body []byte, err error) {
	if b, ok := s.cache.files.Get(name); ok {
		return b.([]byte), nil
	}
	body, err = ioutil.ReadFile(name)
	if err != nil {
		return
	}
	s.cache.files.Set(name, body, cache.DefaultExpiration)
	return
}

// start - Starts watching for file changes in a goroutine.
func (s *Storage) start() (err error) {
	s.events = make(chan notify.EventInfo, 1)
	err = notify.Watch(s.dirname, s.events, notify.All)
	if err != nil {
		return
	}
	go s.watch()
	return
}

// watch - Watches for changes in templates files.
// If change event is emitted, compiled template is deleted from cache.
func (s *Storage) watch() {
	for event := range s.events {
		s.cache.files.Delete(event.Path())
		s.cache.templates.Delete(event.Path())
	}
}
