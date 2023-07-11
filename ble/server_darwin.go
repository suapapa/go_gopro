//go:build darwin

package ble

import (
	"fmt"

	goble "github.com/go-ble/ble"
)

func NewPlatformDeviceWithName(name string, opts ...goble.Option) (goble.Device, error) {
	return nil, fmt.Errorf("not implemented")
}
