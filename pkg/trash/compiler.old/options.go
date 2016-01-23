package compiler

// Option - Compiler option.
type Option func(*options)

type options struct {
	// dirs - List of directories containing components.
	dirs []string

	// watch - When true compiler looks for changes in directories.
	watch bool
}

// WithDirs - Appends list of directories containing components.
// Compiler looks in all directories recursively and when finds
// `component.json` file it reads it and inserts a new component.
func WithDirs(dirs ...string) Option {
	return func(o *options) {
		o.dirs = append(o.dirs, dirs...)
	}
}

// WithWatch - When true compiler looks for changes in directories.
// For now watching is naive meaning it recompiles everything,
// it is only used in local development and compiling should be
// fast enough so this will not be a problem.
func WithWatch(watch bool) Option {
	return func(o *options) {
		o.watch = watch
	}
}
