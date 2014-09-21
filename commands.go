package main

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandInit,
	commandStart,
}

var commandInit = cli.Command{
	Name:  "init",
	Usage: "Create a configuration file",
	Description: `
`,
	Action: doInit,
}

var commandStart = cli.Command{
	Name:  "start",
	Usage: "Start the main program",
	Description: `
`,
	Action: doStart,
	Flags:  flagStart,
}
