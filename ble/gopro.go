package ble

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	goble "github.com/go-ble/ble"
	"github.com/pkg/errors"
)

type GoPro struct {
	cln goble.Client
	p   *goble.Profile
	chs map[Characteristic]*goble.Characteristic
}

func ScanGoPro(opts ...goble.Option) (*GoPro, error) {
	dev, err := newDevice(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new device")
	}

	goble.SetDefaultDevice(dev)

	ctx := goble.WithSigHandler(context.WithTimeout(context.Background(), 10*time.Second))
	filter := func(a goble.Advertisement) bool {
		svcs := a.Services()
		for _, svc := range svcs {
			log.Println(svc)
			if svc.Equal(svcUUIDControlAndQuery) {
				return true
			}
		}
		return false
	}

	cln, err := goble.Connect(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to GoPro")
	}

	p, err := cln.DiscoverProfile(true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to discover profile")
	}

	ret := &GoPro{
		cln: cln,
		p:   p,
		chs: makeCharacteristicMap(p),
	}

	return ret, nil
}

func (g *GoPro) Close() error {
	// Unsubscribe from notifications
	exitC := g.cln.Disconnected()
	err := g.cln.CancelConnection()
	if err != nil {
		return errors.Wrap(err, "failed to cancel connection")
	}
	<-exitC

	return nil
}

func (g *GoPro) KeepAlive() error {
	chrReq, err := g.getChr(Setting)
	if err != nil {
		return errors.Wrap(err, "failed to get chr")
	}
	chrResp, err := g.getChr(SettingResponse)
	if err != nil {
		return errors.Wrap(err, "failed to get chr")
	}

	doneC := make(chan error)
	// response will be sent to GP-0075, chrSettingResponse
	notiHandler := func(req []byte) {
		if bytes.Equal(req, []byte{0x02, 0x5B, 0x00}) {
			doneC <- nil
			return
		}
		doneC <- errors.Errorf("unexpected response: %v", req)
	}
	err = g.cln.Subscribe(chrResp, false, notiHandler)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to chr")
	}
	defer g.cln.Unsubscribe(chrResp, false)

	// send followin payload to GP-0074, chrSetting
	p := []byte{0x03, 0x5b, 0x01, 0x42}
	err = g.cln.WriteCharacteristic(chrReq, p, false)
	if err != nil {
		return errors.Wrap(err, "failed to write chr")
	}

	select {
	case err := <-doneC:
		if err != nil {
			return errors.Wrap(err, "failed to receive notification")
		}
	case <-time.After(5 * time.Second):
		return errors.New("timeout")
	}
	return nil
}

func (g *GoPro) SetShutter(on bool) error {
	chrReq, err := g.getChr(Command)
	if err != nil {
		return errors.Wrap(err, "failed to get chr")
	}
	chrResp, err := g.getChr(CommandResponse)
	if err != nil {
		return errors.Wrap(err, "failed to get chr")
	}

	doneC := make(chan error)
	// response will be sent to GP-0075, chrSettingResponse
	notiHandler := func(req []byte) {
		// TODO: make tlv from command id
		if bytes.Equal(req, []byte{0x02, 0x01, 0x00}) {
			doneC <- nil
			return
		}
		doneC <- errors.Errorf("unexpected response: %v", req)
	}
	err = g.cln.Subscribe(chrResp, false, notiHandler)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to chr")
	}
	defer g.cln.Unsubscribe(chrResp, false)

	var param []byte
	if on {
		param = []byte{0x01}
	} else {
		param = []byte{0x00}
	}
	p, err := makeTlvCmdWithParam(cmdSetShutter, param)
	if err != nil {
		return errors.Wrap(err, "failed to make tlv")
	}
	err = g.cln.WriteCharacteristic(chrReq, p, false)
	if err != nil {
		return errors.Wrap(err, "failed to write chr")
	}

	select {
	case err := <-doneC:
		if err != nil {
			return errors.Wrap(err, "failed to receive notification")
		}
	case <-time.After(5 * time.Second):
		return errors.New("timeout")
	}
	return nil
}

func (g *GoPro) getChr(c Characteristic) (*goble.Characteristic, error) {
	ch, ok := g.chs[c]
	if !ok {
		return nil, fmt.Errorf("chr %d not found", c)
	}
	return ch, nil
}
