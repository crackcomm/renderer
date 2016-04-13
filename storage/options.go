package storage

import "time"

// Option - Storage option setter.
type Option func(*options)

type options struct {
	dirname          string
	cacheExpiration  time.Duration
	cleanupInterval  time.Duration
	removeWhitespace bool
}

func newOptions(opts ...Option) (o *options) {
	o = &options{
		cacheExpiration: 5 * time.Minute,
		cleanupInterval: 1 * time.Minute,
	}
	for _, opt := range opts {
		opt(o)
	}
	return
}

// WithDir - Sets storage directory.
func WithDir(dirname string) Option {
	return func(o *options) {
		o.dirname = dirname
	}
}

// WithCacheExpiration - Sets cache expiration.
func WithCacheExpiration(cacheExpiration time.Duration) Option {
	return func(o *options) {
		o.cacheExpiration = cacheExpiration
	}
}

// WithCacheCleanupInterval - Sets cache cleanup interval.
func WithCacheCleanupInterval(cleanupInterval time.Duration) Option {
	return func(o *options) {
		o.cleanupInterval = cleanupInterval
	}
}

// WithWhitespaceRemoval - Enables total whitespaces removal (only repeated).
func WithWhitespaceRemoval(removeWhitespace bool) Option {
	return func(o *options) {
		o.removeWhitespace = removeWhitespace
	}
}
