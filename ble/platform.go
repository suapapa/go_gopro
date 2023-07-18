package ble

import (
	goble "github.com/go-ble/ble"
)

// func NewDeviceWithName(name string, opts ...goble.Option) (goble.Device, error) {
// 	return newPlatformDeviceWithName(name, opts...)
// }

func newDevice(opts ...goble.Option) (d goble.Device, err error) {
	return newPlatformDevice(opts...)
}
