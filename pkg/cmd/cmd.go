package cmd

import (
	"time"

	"bitbucket.org/moovie/renderer/pkg/renderer"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
)

// Commands - List of renderer commands.
var Commands = []cli.Command{
	compilerCommand,
}

var compilerCommand = cli.Command{
	Name:  "compiler",
	Usage: "compiles components from directories",
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
		dirname := c.String("dir")
		if dirname == "" {
			glog.Fatal("Components directory needs to be set in --dir.")
		}

		storage, err := renderer.NewStorage(dirname, 15*time.Minute, 5*time.Minute)
		if err != nil {
			glog.Fatalf("[storage] %v", err)
		}
		defer storage.Close()

		// compiler, err := renderer.NewCompiler(storage)
		// if err != nil {
		// 	glog.Fatalf("[compiler] %v", err)
		// }

		compiler := renderer.NewCompiler(storage)

		// cmp := compiler.New(
		// 	compiler.WithWatch(c.Bool("watch")),
		// 	compiler.WithDirs(c.StringSlice("dir")...),
		// )

		// if c.Bool("web") {
		// 	glog.Infof("[api] starting on %s", c.String("listen-addr"))
		//
		// 	ctx := compiler.NewContext(context.Background(), cmp)
		// 	go serveAPI(ctx, c)
		// }
	},
}

// func serveAPI(ctx context.Context, c *cli.Context) {
// 	var chain xhandler.Chain
//
// 	// Add close notifier handler so context is cancelled when the client closes
// 	// the connection
// 	chain.UseC(xhandler.CloseHandler)
//
// 	// Add timeout handler
// 	chain.UseC(xhandler.TimeoutHandler(c.Duration("render-timeout")))
//
// 	handler := chain.HandlerCtx(ctx, xhandler.HandlerFuncC(api.Handler))
//
// 	err := http.ListenAndServe(c.String("listen-addr"), handler)
// 	if err != nil {
// 		glog.Fatalf("[api] %v", err)
// 	}
// }
