package cmd

import (
	"bitbucket.org/moovie/renderer/compiler"
	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

var compileCommand = cli.Command{
	Name:  "compile",
	Usage: "compiles components from directories",
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "dir",
			Usage: "directory containing components",
		},
		cli.BoolFlag{
			Name:  "watch",
			Usage: "watch for changes in components",
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

		// Starts compiler and watches (if --watch) for changes.
		if err := cmp.Start(context.Background()); err != nil {
			glog.Fatalf("[compiler] %v", err)
		}
	},
}
