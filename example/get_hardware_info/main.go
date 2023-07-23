package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	gpble "github.com/suapapa/go_gopro/ble"
)

var (
	adtID string
)

func main() {
	flag.StringVar(&adtID, "a", "hci0", "bluetooth adapter")
	flag.Parse()

	log.Println("Scanning GoPro...")
	gps, err := gpble.ScanGoPro(adtID, time.Second*5)
	if err != nil {
		panic(err)
	}

	log.Println("Found GoPro(s):")
	for _, gp := range gps {
		fmt.Println(gp)
	}

	gp := gps[0]
	gp.Connect()
	defer gp.Close()

	go func() {
		tkr := time.NewTicker(time.Second * 3)
		defer tkr.Stop()
		for {
			select {
			case <-tkr.C:
				log.Println("KeepAlive")
				gp.KeepAlive()
			}
		}
	}()

	log.Println("Getting hardware info...")
	info, err := gp.GetHardwareInfo()
	if err != nil {
		panic(err)
	}

	fmt.Println(info)
	log.Println("Done")
}
