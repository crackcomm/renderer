package watcher

import (
	"fmt"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

// Watcher - Watches for changes and flushes caches.
type Watcher struct {
	flusher CacheFlusher
	events  chan notify.EventInfo
}

// CacheFlusher - Cache flusher interface.
type CacheFlusher interface {
	FlushCache()
}

// Start - Creates a new watcher which flushes caches.
// Starts it in a separate goroutine.
func Start(path string, cache CacheFlusher) (w *Watcher, err error) {
	w, err = New(path, cache)
	if err != nil {
		return
	}
	go w.Start()
	return
}

// New - Creates a new watcher which flushes caches.
func New(path string, cache CacheFlusher) (*Watcher, error) {
	events := make(chan notify.EventInfo, 1)
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, "...")
	err = notify.Watch(path, events, notify.All)
	if err != nil {
		return nil, fmt.Errorf("watching on %q: %v", path, err)
	}
	return &Watcher{
		events:  events,
		flusher: cache,
	}, nil
}

// Start - Starts watching for file changes in current goroutine.
func (w *Watcher) Start() {
	for range w.events {
		w.flusher.FlushCache()
	}
}

// Stop - Stops watcher.
func (w *Watcher) Stop() {
	if w.events != nil {
		notify.Stop(w.events)
		close(w.events)
	}
}
