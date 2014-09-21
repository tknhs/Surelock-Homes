package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

func doInit(c *cli.Context) {
	if err := setConfig(); err != nil {
		panic(err)
	}
}

func doStart(c *cli.Context) {
	fmt.Println("start")

	config, err := getConfig()
	if err != nil {
		log.Fatalf("at first, \"$ surelock-homes init\"\n", err)
	}
	fmt.Println(c.String("serial"))
	fmt.Println(config)
}
