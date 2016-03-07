package renderer

import "time"

// StorageOption - Storage option setter.
type StorageOption func(*storageOptions)

type storageOptions struct {
	dirname          string
	cacheExpiration  time.Duration
	cleanupInterval  time.Duration
	removeWhitespace bool
	watchingChanges  bool
}

// WithWatching - Enables storage watching for changes.
func WithWatching(enable ...bool) StorageOption {
	return func(o *storageOptions) {
		if len(enable) >= 1 {
			o.watchingChanges = enable[0]
		} else {
			o.watchingChanges = true
		}
	}
}

// WithDir - Sets storage directory.
func WithDir(dirname string) StorageOption {
	return func(o *storageOptions) {
		o.dirname = dirname
	}
}

// WithCacheExpiration - Sets cache expiration.
func WithCacheExpiration(cacheExpiration time.Duration) StorageOption {
	return func(o *storageOptions) {
		o.cacheExpiration = cacheExpiration
	}
}

// WithCacheCleanupInterval - Sets cache cleanup interval.
func WithCacheCleanupInterval(cleanupInterval time.Duration) StorageOption {
	return func(o *storageOptions) {
		o.cleanupInterval = cleanupInterval
	}
}

// WithWhitespaceRemoval - Enables total whitespaces removal (only repeated).
func WithWhitespaceRemoval(removeWhitespace bool) StorageOption {
	return func(o *storageOptions) {
		o.removeWhitespace = removeWhitespace
	}
}
