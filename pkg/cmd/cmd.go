package cmd

import (
	"net/http"
	"time"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"golang.org/x/net/context"

	"bitbucket.org/moovie/renderer/pkg/renderer"
	"bitbucket.org/moovie/renderer/pkg/web"
)

// Commands - List of renderer commands.
var Commands = []cli.Command{
	webCommand,
}

var webCommand = cli.Command{
	Name:  "web",
	Usage: "starts renderer web API",
	Flags: []cli.Flag{
		// Compiler options
		cli.StringFlag{
			Name:  "dir",
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
		// Get components directory from --dir flag
		// Print fatal error if not set
		if c.String("dir") == "" {
			glog.Fatal("Components directory needs to be set in --dir.")
		}

		// Create a new storage in directory from --dir flag
		storage, err := renderer.NewStorage(
			renderer.WithDir(c.String("dir")),
			renderer.WithCacheExpiration(c.Duration("cache-expiration")),
			renderer.WithCacheCleanupInterval(c.Duration("cache-cleanup")),
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
		ctx = renderer.CompilerCtx(ctx, compiler)

		// Create a web server http handler
		api := web.NewAPI(
			web.WithContext(ctx),
		)

		glog.Infof("[api] starting on %s", c.String("listen-addr"))

		// Start http server
		err = http.ListenAndServe(c.String("listen-addr"), api)
		if err != nil {
			glog.Fatalf("[api] %v", err)
		}
	},
}
