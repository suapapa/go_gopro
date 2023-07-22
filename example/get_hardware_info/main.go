package main

import (
	"flag"
	"fmt"
	"time"

	gpble "github.com/suapapa/go_gopro/ble"
)

var (
	adtID string
)

func main() {
	flag.StringVar(&adtID, "a", "hci0", "bluetooth adapter")
	flag.Parse()

	gps, err := gpble.ScanGoPro(adtID, time.Second*5)
	if err != nil {
		panic(err)
	}

	for _, gp := range gps {
		fmt.Println(gp)
	}

	gps[0].Connect()
}
