package cmd

import (
	"net/http"
	"time"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"

	"bitbucket.org/moovie/renderer/api"
	"bitbucket.org/moovie/renderer/compiler"
)

var compilerCommand = cli.Command{
	Name:  "compiler",
	Usage: "compiles components from directories",
	Flags: []cli.Flag{
		// Compiler options
		cli.StringSliceFlag{
			Name:  "dir",
			Usage: "directory containing components",
		},
		cli.BoolFlag{
			Name:  "watch",
			Usage: "watch for changes in components",
		},

		// Web interface options
		cli.BoolFlag{
			Name:  "web",
			Usage: "enables web interface",
		},
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
		dirs := c.StringSlice("dir")
		if len(dirs) == 0 {
			glog.Fatal("No directories were set. At least one --dir is required.")
		}

		cmp := compiler.New(
			compiler.WithWatch(c.Bool("watch")),
			compiler.WithDirs(c.StringSlice("dir")...),
		)

		if c.Bool("web") {
			glog.Infof("[api] starting on %s", c.String("listen-addr"))

			ctx := compiler.NewContext(context.Background(), cmp)
			go serveAPI(ctx, c)
		}

		// Starts compiler and watches (if --watch) for changes.
		if err := cmp.Start(context.Background()); err != nil {
			glog.Fatalf("[compiler] %v", err)
		}

		if c.Bool("web") {
			select {}
		}
	},
}

func serveAPI(ctx context.Context, c *cli.Context) {
	var chain xhandler.Chain

	// Add close notifier handler so context is cancelled when the client closes
	// the connection
	chain.UseC(xhandler.CloseHandler)

	// Add timeout handler
	chain.UseC(xhandler.TimeoutHandler(c.Duration("render-timeout")))

	handler := chain.HandlerCtx(ctx, xhandler.HandlerFuncC(api.Handler))

	err := http.ListenAndServe(c.String("listen-addr"), handler)
	if err != nil {
		glog.Fatalf("[api] %v", err)
	}
}
