//go:build linux

package ble

import (
	goble "github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func NewPlatformDeviceWithName(name string, opts ...goble.Option) (goble.Device, error) {
	return linux.NewDeviceWithNameAndHandler(name, nil, opts...)
}

type NotifyHandler func(goble.Request, goble.Notifier)

func (f NotifyHandler) Handle(req goble.Request, n goble.Notifier) {
	f(req, n)
}
