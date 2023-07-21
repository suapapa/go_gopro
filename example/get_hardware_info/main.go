package main

import (
	"fmt"

	"github.com/go-ble/ble"
	gpble "github.com/suapapa/go_gopro/ble"
)

func main() {
	gp, err := gpble.ScanGoPro(
		ble.OptDeviceID(0),
	)
	if err != nil {
		panic(err)
	}
	defer gp.Close()

	fmt.Printf("GoPro found: %s\n", gp)

	err = gp.KeepAlive()
	if err != nil {
		panic(err)
	}

	hwInfo, err := gp.GetHardwareInfo()
	if err != nil {
		panic(err)
	}
	println(hwInfo)
}
