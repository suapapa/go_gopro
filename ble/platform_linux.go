//go:build linux

package ble

import (
	goble "github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

// func newPlatformDeviceWithName(name string, opts ...goble.Option) (goble.Device, error) {
// 	return linux.NewDeviceWithNameAndHandler(name, nil, opts...)
// }

/*
type NotifyHandler func(goble.Request, goble.Notifier)

func (f NotifyHandler) Handle(req goble.Request, n goble.Notifier) {
	f(req, n)
}
*/

func newPlatformDevice(opts ...goble.Option) (d goble.Device, err error) {
	return linux.NewDevice(opts...)
}
