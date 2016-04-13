package command

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/compiler"
	"github.com/crackcomm/renderer/renderweb"
	"github.com/crackcomm/renderer/storage"
	"github.com/crackcomm/renderer/watcher"

	// Profiler
	_ "net/http/pprof"
)

// Commands - List of renderer commands.
var Commands = []cli.Command{
	Web,
}

// Web - Web command.
var Web = cli.Command{
	Name:  "web",
	Usage: "starts renderer web API",
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
			Usage: "web interface listening address",
			Value: "127.0.0.1:6660",
		},
		cli.DurationFlag{
			Name:  "render-timeout",
			Usage: "component render timeout",
			Value: 5 * time.Second,
		},

		// HTTP server flags
		cli.DurationFlag{
			Name:   "http-read-timeout",
			EnvVar: "HTTP_READ_TIMEOUT",
			Usage:  "http server read timeout",
			Value:  time.Minute,
		},
		cli.DurationFlag{
			Name:   "http-write-timeout",
			EnvVar: "HTTP_WRITE_TIMEOUT",
			Usage:  "http server write timeout",
			Value:  time.Minute,
		},

		// Profiler
		cli.StringFlag{
			Name:  "pprof-addr",
			Usage: "pprof listening address",
		},
	},
	Action: func(c *cli.Context) {
		// Get components directory from --components flag
		// Print fatal error if not set
		if c.String("components") == "" {
			glog.Fatal("--components flag cannot be empty")
		}

		// Create a new storage in directory from --components flag
		storage, err := storage.New(
			storage.WithDir(c.String("components")),
			storage.WithCacheExpiration(c.Duration("cache-expiration")),
			storage.WithCacheCleanupInterval(c.Duration("cache-cleanup")),
			storage.WithWhitespaceRemoval(c.Bool("compress")),
		)
		if err != nil {
			glog.Fatalf("[storage] %v", err)
		}
		defer storage.Close()

		// Create a compiler from storage
		comp := compiler.New(storage, c.Duration("cache-expiration"), c.Duration("cache-cleanup"))

		// Create a context with compiler
		ctx := compiler.NewContext(context.Background(), comp)

		// Turn routes into HTTP handler
		api, err := constructHandler(c.StringSlice("routes")...)
		if err != nil {
			glog.Fatalf("[routes] %v", err)
		}

		// Construct API handler
		handler := &atomicHandler{
			Context:  ctx,
			Current:  xhandler.New(ctx, api),
			Watching: c.Bool("watch"),
			Routes:   c.StringSlice("routes"),
			Mutex:    new(sync.RWMutex),
		}

		if c.Bool("watch") {
			// Start watching for changes in components directory
			var w *watcher.Watcher
			w, err = watcher.Start(c.String("components"), comp)
			if err != nil {
				glog.Fatal(err)
			}
			defer w.Stop()

			// Start watching for changes in routes
			for _, filename := range c.StringSlice("routes") {
				var watch *watcher.Watcher
				watch, err = watcher.Start(filename, handler)
				if err != nil {
					glog.Fatal(err)
				}
				defer watch.Stop()
			}
		}

		// Start profiler if enabled
		if pprofaddr := c.String("pprof-addr"); pprofaddr != "" {
			go func() {
				glog.Infof("[pprof] starting listener on %s", pprofaddr)
				if err := http.ListenAndServe(pprofaddr, nil); err != nil {
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

		glog.Infof("[server] starting listener on %s", c.String("listen-addr"))
		if err = server.ListenAndServe(); err != nil {
			glog.Fatalf("[server] %v", err)
		}
	},
}

type atomicHandler struct {
	Context context.Context
	Current http.Handler

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
	h, err := constructHandler(handler.Routes...)
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

func constructHandler(filenames ...string) (_ xhandler.HandlerC, err error) {
	if len(filenames) == 0 {
		return renderweb.New(), nil
	}

	routes, err := constructRoutes(filenames...)
	if err != nil {
		return
	}

	// Turn routes into HTTP handler
	return routes.Construct()
}

// constructRoutes - Constructs routes map from multiple filenames.
func constructRoutes(filenames ...string) (res renderweb.Routes, err error) {
	res = make(renderweb.Routes)
	for _, filename := range filenames {
		var routes renderweb.Routes
		routes, err = renderweb.RoutesFromFile(filename)
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
