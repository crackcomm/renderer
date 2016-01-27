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
		dirname := c.String("dir")
		if dirname == "" {
			glog.Fatal("Components directory needs to be set in --dir.")
		}

		// Create a new storage in directory from --dir flag
		storage, err := renderer.NewStorage(dirname, 15*time.Minute, 5*time.Minute)
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
		w := web.NewAPI(
			web.WithCtx(ctx),
		)

		glog.Infof("[api] starting on %s", c.String("listen-addr"))

		// Start http server
		err = http.ListenAndServe(c.String("listen-addr"), w)
		if err != nil {
			glog.Fatalf("[api] %v", err)
		}
	},
}
