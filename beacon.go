package main

import (
	"fmt"
	"log"
	"os/exec"
)

func BeaconScan() {
	beacon, err := exec.Command("node", "path/to/reciever.js").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string([]byte(beacon)))
}
