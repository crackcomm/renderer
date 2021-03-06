package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"

	"tower.pro/renderer/command"
)

var version string

func main() {
	defer glog.Flush()
	flag.CommandLine.Parse([]string{"-logtostderr"})

	app := cli.NewApp()
	app.Name = "renderer"
	app.Usage = "components compiler, renderer and command line tool with web interface"
	app.Version = version
	app.HideVersion = true
	app.Commands = command.Commands
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "v",
			Usage:  "verbosity (disables logs if < 0)",
			EnvVar: "VERBOSITY",
		},
	}
	app.Before = glogFlags
	app.Run(os.Args)
}

func glogFlags(c *cli.Context) error {
	v := c.Int("v")
	if v < 0 {
		flag.Set("logtostderr", "false")
	} else if v > 0 {
		flag.Set("v", fmt.Sprintf("%d", v))
	}
	return nil
}
