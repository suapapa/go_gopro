package ble

import goble "github.com/go-ble/ble"

func newDevice(opts ...goble.Option) (d goble.Device, err error) {
	return newPlatformDevice(opts...)
}
