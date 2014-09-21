package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "surelock-homes"
	app.Version = Version
	app.Usage = "Home Lock Management"
	app.Commands = Commands

	app.Run(os.Args)
}
