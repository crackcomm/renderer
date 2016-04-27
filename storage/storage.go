package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/patrickmn/go-cache"
	"github.com/rjeczalik/notify"

	"tower.pro/renderer/components"
	"tower.pro/renderer/template"
)

// New - Creates new components storage.
func New(opts ...Option) (s *Storage, err error) {
	o := newOptions(opts...)
	o.dirname, err = filepath.Abs(o.dirname)
	if err != nil {
		return
	}
	return &Storage{
		opts: o,
		cache: &storageCache{
			components: cache.New(o.cacheExpiration, o.cleanupInterval),
			templates:  cache.New(o.cacheExpiration, o.cleanupInterval),
			files:      cache.New(o.cacheExpiration, o.cleanupInterval),
		},
	}, nil
}

// Storage - Components storage.
type Storage struct {
	opts *options

	events chan notify.EventInfo
	cache  *storageCache
}

type storageCache struct {
	components *cache.Cache
	templates  *cache.Cache
	files      *cache.Cache
}

// Text - Returns file content as Template interface.
func (s *Storage) Text(path string) (t template.Template, err error) {
	path = filepath.Join(s.opts.dirname, path)
	body, err := s.read(path)
	if err != nil {
		return
	}
	return template.TextBytes(body), nil
}

// Template - Compiles template by file path and saves in cache.
// Returns cached template if already compiled and not changed.
func (s *Storage) Template(path string) (t template.Template, err error) {
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
func (s *Storage) Component(name string) (c *components.Component, err error) {
	path := strings.Replace(name, ".", string(os.PathSeparator), -1)
	path = filepath.Join(s.opts.dirname, path, "component.yaml")
	if tmp, ok := s.cache.components.Get(path); ok {
		return tmp.(*components.Component), nil
	}
	body, err := s.read(path)
	if err != nil {
		return
	}
	c = new(components.Component)
	err = yaml.Unmarshal(body, c)
	if err != nil {
		return
	}
	if c.Name == "" {
		c.Name = name
	}
	s.cache.components.Set(path, c, cache.DefaultExpiration)
	return
}

// Close - Destroys caches and stops watching for changes.
func (s *Storage) Close() (err error) {
	s.FlushCache()
	return
}

// read - reads file content or returns cached byte array
func (s *Storage) read(path string) (body []byte, err error) {
	if b, ok := s.cache.files.Get(path); ok {
		return b.([]byte), nil
	}
	body, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	if s.opts.removeWhitespace {
		body = CleanWhitespaces(body)
	}
	s.cache.files.Set(path, body, cache.DefaultExpiration)
	return
}

// FlushCache - Flushes storage cache.
func (s *Storage) FlushCache() {
	s.cache.files.Flush()
	s.cache.templates.Flush()
	s.cache.components.Flush()
}
