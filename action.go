package main

import (
	"log"

	"github.com/codegangsta/cli"
)

func doInit(c *cli.Context) {
	if err := SetConfig(); err != nil {
		log.Fatalf("failed to create a setting file\n", err)
	}
}

func doStart(c *cli.Context) {
	// read a configuration file
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("at first, \"$ surelock-homes init\"\n", err)
	}

	// command option check
	cmdOptionSerial := c.String("serial")
	if cmdOptionSerial != "SERIAL PORT" {
		config.SerialPort.Serial = cmdOptionSerial
	}

	// twitter initial
	token := TwitterInit(config.Twitter)
	// serial initial
	serialObject, err := SerialInit(config.SerialPort)
	if err != nil {
		log.Fatalf("failed to open the serial port\n", err)
	}

	log.Println("start TwitterStreaming and BeaconScanning...")

	for {
		timestamp := make(chan string)
		go TwitterStreaming(timestamp, token, config.Twitter.Account)
		go BeaconScan(timestamp, config.Bluetooth)

		ts1 := <-timestamp
		ts2 := <-timestamp

		// door doesn't open when the difference exceeds the 5 minutes
		if timediff := TimeDiff(ts1, ts2); timediff >= 300 {
			continue
		}

		// send a open command
		err = SerialWrite(serialObject, "OC0")
		if err != nil {
			log.Fatalf("failed to write\n", err)
		}

		// twitter post
		err = TwitterPost(token, "[from Surelock-Homes] The door has opened.")
		if err != nil {
			log.Fatalf("failed to post a tweet\n", err)
		}
	}
}
