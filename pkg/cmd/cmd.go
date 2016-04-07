package cmd

import (
	"net/http"
	"time"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/pkg/renderer"
	"github.com/crackcomm/renderer/pkg/routes"
	"github.com/crackcomm/renderer/pkg/web"
)

// Commands - List of renderer commands.
var Commands = []cli.Command{
	CommandWeb,
}

// CommandWeb - Web command.
var CommandWeb = cli.Command{
	Name:  "web",
	Usage: "starts renderer web API",
	Flags: []cli.Flag{
		// Compiler options
		cli.StringFlag{
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
	},
	Action: func(c *cli.Context) {
		// Get components directory from --components flag
		// Print fatal error if not set
		if c.String("components") == "" {
			glog.Fatal("--components flag cannot be empty")
		}

		// Create a new storage in directory from --components flag
		storage, err := renderer.NewStorage(
			renderer.WithDir(c.String("components")),
			renderer.WithCacheExpiration(c.Duration("cache-expiration")),
			renderer.WithCacheCleanupInterval(c.Duration("cache-cleanup")),
			renderer.WithWatching(c.Bool("watch")),
			renderer.WithWhitespaceRemoval(c.Bool("compress")),
		)
		if err != nil {
			glog.Fatalf("[storage] %v", err)
		}
		defer storage.Close()

		// Create a compiler from storage
		compiler := renderer.NewCompiler(storage)

		// Base context
		ctx := context.Background()

		// Create a context with compiler
		ctx = renderer.WithCompiler(ctx, compiler)

		glog.Infof("[api] starting on %s", c.String("listen-addr"))

		// Start serving routes if set
		var api xhandler.HandlerC
		if fname := c.String("routes"); fname != "" {
			// Read routes from file
			r, err := routes.FromFile(fname)
			if err != nil {
				glog.Fatalf("[routes] %v", err)
			}

			// Turn routes into HTTP handler
			api, err = r.Construct()
			if err != nil {
				glog.Fatalf("[routes] %v", err)
			}
		} else {
			api = web.New()
		}

		// Start http server
		server := &http.Server{
			Addr:           c.String("listen-addr"),
			Handler:        xhandler.New(ctx, api),
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 64 * 1024,
		}

		if err = server.ListenAndServe(); err != nil {
			glog.Fatalf("[server] %v", err)
		}
	},
}
