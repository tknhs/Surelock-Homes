package main

import (
	"github.com/codegangsta/cli"
)

// Flag
var flagStart = []cli.Flag{
	startSerial,
}
var startSerial = cli.StringFlag{
	Name:  "serial",
	Value: "SERIAL PORT",
	Usage: "Setting the serial port",
}
