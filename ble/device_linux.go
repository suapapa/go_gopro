package ble

import (
	goble "github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func newPlatformDevice(opts ...goble.Option) (d goble.Device, err error) {
	return linux.NewDevice(opts...)
}
