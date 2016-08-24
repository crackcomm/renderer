package command

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/rs/xhandler"

	"tower.pro/renderer/compiler"
	"tower.pro/renderer/renderer"
	"tower.pro/renderer/storage"
	"tower.pro/renderer/watcher"

	// Profiler
	_ "net/http/pprof"
)

// DefaultWebOptions - Default renderer options.
// You can append options here like renderer.WithMiddleware.
var DefaultWebOptions = []renderer.Option{
	// It means we won't spit out component JSON, just HTML
	// in case of API routes like one using json.encode they
	// will still work as expected with no difference
	renderer.WithAlwaysHTML(true),
}

// Commands - List of renderer commands.
var Commands = []cli.Command{
	Web,
}

// Web - Web command.
var Web = cli.Command{
	Name:  "server",
	Usage: "renderer server",
	Flags: []cli.Flag{
		// Compiler options
		cli.StringSliceFlag{
			Name:  "routes",
			Usage: "file containing routes in yaml format",
		},
		cli.StringFlag{
			Name:  "components",
			Usage: "directory containing components",
		},
		cli.BoolFlag{
			Name:  "watch",
			Usage: "watch for changes in components",
		},
		cli.BoolFlag{
			Name:  "compress",
			Usage: "removes repeated whitespaces",
		},
		cli.DurationFlag{
			Name:  "cache-expiration",
			Usage: "cache expiration time",
			Value: 15 * time.Minute,
		},
		cli.DurationFlag{
			Name:  "cache-cleanup",
			Usage: "cache cleanup interval",
			Value: 5 * time.Minute,
		},

		// Web server options
		cli.StringFlag{
			Name:  "listen-addr",
			Usage: "renderer interface listening address",
			Value: "127.0.0.1:6660",
		},
		cli.DurationFlag{
			Name:  "render-timeout",
			Usage: "component render timeout",
			Value: 5 * time.Second,
		},

		// HTTP server flags
		cli.DurationFlag{
			Name:   "renderer-read-timeout",
			EnvVar: "RENDERER_READ_TIMEOUT",
			Usage:  "renderer server read timeout",
			Value:  time.Minute,
		},
		cli.DurationFlag{
			Name:   "renderer-write-timeout",
			EnvVar: "RENDERER_WRITE_TIMEOUT",
			Usage:  "renderer server write timeout",
			Value:  time.Minute,
		},

		// Tracing and profiling
		cli.StringFlag{
			Name:   "debug-addr",
			EnvVar: "DEBUG_ADDR",
			Usage:  "debug listening address",
		},
		cli.BoolFlag{
			Name:   "tracing",
			EnvVar: "TRACING",
			Usage:  "enable tracing (use with --debug-addr)",
		},
	},
	Action: func(c *cli.Context) (err error) {
		// Get components directory from --components flag
		// Print fatal error if not set
		if c.String("components") == "" {
			return errors.New("--components flag cannot be empty")
		}

		// Create a new storage in directory from --components flag
		storage, err := storage.New(
			storage.WithDir(c.String("components")),
			storage.WithCacheExpiration(c.Duration("cache-expiration")),
			storage.WithCacheCleanupInterval(c.Duration("cache-cleanup")),
			storage.WithWhitespaceRemoval(c.Bool("compress")),
		)
		if err != nil {
			return fmt.Errorf("[storage] %v", err)
		}
		defer storage.Close()

		// Create a compiler from storage
		comp := compiler.New(storage)

		// Create a context with compiler
		ctx := compiler.NewContext(context.Background(), comp)

		if c.Bool("tracing") {
			DefaultWebOptions = append(DefaultWebOptions, renderer.WithTracing())
		}

		// Turn routes into HTTP handler
		api, err := constructHandler(c.StringSlice("routes"), DefaultWebOptions)
		if err != nil {
			return fmt.Errorf("[routes] %v", err)
		}

		// Construct API handler
		handler := &atomicHandler{
			Context:  ctx,
			Current:  xhandler.New(ctx, api),
			Options:  DefaultWebOptions,
			Watching: c.Bool("watch"),
			Routes:   c.StringSlice("routes"),
			Mutex:    new(sync.RWMutex),
		}

		if c.Bool("watch") {
			// Start watching for changes in components directory
			var w *watcher.Watcher
			w, err = watcher.Start(c.String("components"), storage)
			if err != nil {
				return
			}
			defer w.Stop()

			// Start watching for changes in routes
			for _, filename := range c.StringSlice("routes") {
				var watch *watcher.Watcher
				watch, err = watcher.Start(filename, handler)
				if err != nil {
					return
				}
				defer watch.Stop()
			}
		}

		// Start profiler if enabled
		if addr := c.String("pprof-addr"); addr != "" {
			go func() {
				if err = debugServer(addr); err != nil {
					glog.Fatal(err)
				}
			}()
		}

		// Construct http server
		server := &http.Server{
			Addr:           c.String("listen-addr"),
			Handler:        handler,
			ReadTimeout:    c.Duration("http-read-timeout"),
			WriteTimeout:   c.Duration("http-write-timeout"),
			MaxHeaderBytes: 64 * 1024,
		}

		glog.Infof("[renderer] starting server on %s", c.String("listen-addr"))
		return server.ListenAndServe()
	},
}

type atomicHandler struct {
	Context context.Context
	Current http.Handler
	Options []renderer.Option

	Mutex    *sync.RWMutex
	Watching bool
	Routes   []string
}

// FlushCache - Flushes routes cache. Reads them and constructs handler.
func (handler *atomicHandler) FlushCache() {
	// Construct handler from routes
	h, err := handler.construct()
	if err != nil {
		glog.Fatalf("[routes] %v", err)
	}

	// Lock mutex and exchange handler
	handler.Mutex.Lock()
	handler.Current = h
	handler.Mutex.Unlock()
}

func (handler *atomicHandler) construct() (_ http.Handler, err error) {
	h, err := constructHandler(handler.Routes, handler.Options)
	if err != nil {
		return
	}
	return xhandler.New(handler.Context, h), nil
}

func (handler *atomicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If not watching don't use mutexe
	if !handler.Watching {
		handler.Current.ServeHTTP(w, r)
		return
	}

	// Lock for read and get handler
	handler.Mutex.RLock()
	h := handler.Current
	handler.Mutex.RUnlock()

	// Serve request
	h.ServeHTTP(w, r)
}

func constructHandler(filenames []string, options []renderer.Option) (_ xhandler.HandlerC, err error) {
	if len(filenames) == 0 {
		return renderer.New(), nil
	}

	routes, err := constructRoutes(filenames, options)
	if err != nil {
		return
	}

	// Turn routes into HTTP handler
	return routes.Construct(options...)
}

// constructRoutes - Constructs routes map from multiple filenames.
func constructRoutes(filenames []string, options []renderer.Option) (res renderer.Routes, err error) {
	res = make(renderer.Routes)
	for _, filename := range filenames {
		var routes renderer.Routes
		routes, err = renderer.RoutesFromFile(filename)
		if err != nil {
			return
		}
		for route, handler := range routes {
			if _, exists := res[route]; exists {
				return nil, fmt.Errorf("route %q in %q is not unique", route, filename)
			}
			res[route] = handler
		}
	}
	return
}

func debugServer(addr string) (err error) {
	glog.Infof("[debug] starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
