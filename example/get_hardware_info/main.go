package main

import (
	"github.com/go-ble/ble"
	gpble "github.com/suapapa/go_gopro/ble"
)

func main() {
	gp, err := gpble.ScanGoPro(
		ble.OptCentralRole(),
	)
	if err != nil {
		panic(err)
	}
	defer gp.Close()

	hwInfo, err := gp.GetHardwareInfo()
	if err != nil {
		panic(err)
	}
	println(hwInfo)
}
