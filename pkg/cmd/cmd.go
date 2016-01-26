package cmd

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/rs/xhandler"

	"bitbucket.org/moovie/renderer/pkg/api"
	"bitbucket.org/moovie/renderer/pkg/renderer"
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

		// Web interface options
		cli.StringFlag{
			Name:  "listen-addr",
			Usage: "web interface listening address",
			Value: ":5055",
		},
		cli.DurationFlag{
			Name:  "render-timeout",
			Usage: "component render timeout",
			Value: 5 * time.Second,
		},
	},
	Action: func(c *cli.Context) {
		dirname := c.String("dir")
		if dirname == "" {
			glog.Fatal("Components directory needs to be set in --dir.")
		}

		storage, err := renderer.NewStorage(dirname, 15*time.Minute, 5*time.Minute)
		if err != nil {
			glog.Fatalf("[storage] %v", err)
		}
		defer storage.Close()

		compiler := renderer.NewCompiler(storage)

		var chain xhandler.Chain

		// Add close notifier handler so context is cancelled when the client closes
		// the connection
		chain.UseC(xhandler.CloseHandler)

		// Add timeout handler
		chain.UseC(xhandler.TimeoutHandler(c.Duration("render-timeout")))

		ctx := renderer.NewContext(context.Background(), compiler)
		handler := chain.HandlerCtx(ctx, xhandler.HandlerFuncC(api.Handler))

		glog.Infof("[api] starting on %s", c.String("listen-addr"))
		err = http.ListenAndServe(c.String("listen-addr"), handler)
		if err != nil {
			glog.Fatalf("[api] %v", err)
		}
	},
}
