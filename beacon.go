package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func BeaconScan(beaconTimestamp chan string, beaconConfig BluetoothConfig) {
	isBeacon := true

	for isBeacon {
		beaconInfo, err := exec.Command("node", "path/to/reciever.js").Output()
		if err != nil {
			log.Fatal(err)
		}

		var bc BluetoothConfig
		err = json.Unmarshal(beaconInfo, &bc)
		if err != nil {
			log.Fatalf("can't decode json\n", err)
		}

		if beaconConfig == bc {
			isBeacon = false
		}
	}

	t := time.Now().Unix()
	beaconTimestamp <- strconv.Itoa(int(t))
}
