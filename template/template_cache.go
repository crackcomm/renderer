package template

import (
	"io/ioutil"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rjeczalik/notify"
)

// CachedSet - Cached templates set interface.
type CachedSet struct {
	dirname string
	events  chan notify.EventInfo
	cache   struct {
		files     *cache.Cache
		templates *cache.Cache
	}
}

// NewCachedSet - Creates new cached template set.
// If second parameter is empty uses DefaultCache.
func NewCachedSet(dirname string, defaultExpiration time.Duration, cleanupInterval time.Duration) (s *CachedSet, err error) {
	s = &CachedSet{dirname: dirname}
	s.cache.files = cache.New(defaultExpiration, cleanupInterval)
	s.cache.templates = cache.New(defaultExpiration, cleanupInterval)
	err = s.start()
	return
}

// Text - Returns file content as Template interface.
func (s *CachedSet) Text(name string) (t Template, err error) {
	body, err := s.read(name)
	if err != nil {
		return
	}
	return TextBytes(body), nil
}

// Template - Compiles template by filename and saves in cache.
// Returns cached template if already compiled and not changed.
func (s *CachedSet) Template(name string) (t Template, err error) {
	if tmp, ok := s.cache.templates.Get(name); ok {
		if t, ok := tmp.(Template); ok {
			return t, nil
		}
	}
	body, err := s.read(name)
	if err != nil {
		return
	}
	t, err = FromBytes(body)
	if err != nil {
		return
	}
	s.cache.templates.Set(name, t, cache.DefaultExpiration)
	return
}

// Close - Closes cached set, destroys cache and stops watching for changes.
func (s *CachedSet) Close() (err error) {
	if s.events != nil {
		close(s.events)
	}
	return
}

// read - reads file content
func (s *CachedSet) read(name string) (body []byte, err error) {
	if b, ok := s.cache.files.Get(name); ok {
		if body, ok := b.([]byte); ok {
			return body, nil
		}
	}
	body, err = ioutil.ReadFile(name)
	if err != nil {
		return
	}
	s.cache.files.Set(name, body, cache.DefaultExpiration)
	return
}

// start - Starts watching for file changes in a goroutine.
func (s *CachedSet) start() (err error) {
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
func (s *CachedSet) watch() {
	for event := range s.events {
		s.cache.files.Delete(event.Path())
		s.cache.templates.Delete(event.Path())
	}
}
