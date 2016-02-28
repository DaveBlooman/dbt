package main

import (
	"os"

	"github.com/DaveBlooman/dbt/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "DaveBlooman"
	app.Email = "david.blooman@gmail.com"
	app.Usage = "CLI tool for building RPMs"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
