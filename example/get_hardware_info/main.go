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

	gps, err := gpble.ScanGoPro(adtID, time.Second*10)
	if err != nil {
		panic(err)
	}

	for _, gp := range gps {
		fmt.Println(gp)
	}
}
