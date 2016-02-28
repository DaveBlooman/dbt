package main

import (
	"fmt"
	"os"

	"github.com/DaveBlooman/dbt/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/DaveBlooman/dbt/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{

	{
		Name:   "build",
		Usage:  "Build RPM",
		Action: command.CmdBuild,
		Flags:  []cli.Flag{},
	},

	{
		Name:   "pull",
		Usage:  "Pull Docker image",
		Action: command.CmdPull,
		Flags:  []cli.Flag{},
	},

	{
		Name:   "test",
		Usage:  "",
		Action: command.CmdTest,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
